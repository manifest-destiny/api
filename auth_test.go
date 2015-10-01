package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/test.rsa"     // openssl genrsa -out app.rsa 1024
	pubKeyPath  = "keys/test.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub // (json-ified)
)

var TokenTestData = []struct {
	name              string
	iss, aud, kid     string // token properties
	exp               int64  // expiry
	issuer, apiClient string // auth settings
	valid             bool   // is token valid
}{
	{
		"valid token",
		"accounts.google.com",
		"webapp_client_id",
		"test_key",
		time.Now().Add(time.Hour * 1).Unix(),
		"accounts.google.com",
		"webapp_client_id",
		true,
	},
	{
		"expired time",
		"accounts.google.com",
		"webapp_client_id",
		"test_key",
		time.Now().Add(time.Hour * -1).Unix(),
		"accounts.google.com",
		"webapp_client_id",
		false,
	},
	{
		"invalid api client",
		"accounts.google.com",
		"webapp_client_id",
		"test_key",
		time.Now().Add(time.Hour * 1).Unix(),
		"account.google.com",
		"ios_client_id",
		false,
	},
	{
		"bogus verification key",
		"accounts.google.com",
		"webapp_client_id",
		"bogus_test_key",
		time.Now().Add(time.Hour * 1).Unix(),
		"account.google.com",
		"webapp_client_id",
		false,
	},
	{
		"invalid token issuer",
		"accounts.facebook.com",
		"webapp_client_id",
		"bogus_test_key",
		time.Now().Add(time.Hour * 1).Unix(),
		"account.google.com",
		"webapp_client_id",
		false,
	},
}

func TestTokenAuthorisation(t *testing.T) {
	signBytes, _ := ioutil.ReadFile(privKeyPath)
	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	// Create the token
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	for _, data := range TokenTestData {
		token.Claims["iss"] = data.iss
		token.Claims["aud"] = data.aud
		token.Claims["exp"] = data.exp
		token.Header["kid"] = data.kid

		// get certs
		f, _ := ioutil.ReadFile(pubKeyPath)
		var certs map[string]string
		json.Unmarshal(f, &certs)

		// setup google auth
		a := NewGoogleAuth(certs, data.issuer, data.apiClient)

		// retrieve the token
		tokenStr, _ := token.SignedString(signKey)
		gToken, _ := a.Verify(tokenStr)

		if gToken.Valid != data.valid {
			if data.valid {
				t.Errorf("Expected token to be valid for \"%s\"", data.name)
			} else {
				t.Errorf("Expected token to be invalid for \"%s\"", data.name)
			}
		}
	}
}

func ExamplePlayerData() {
	signBytes, _ := ioutil.ReadFile(privKeyPath)
	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	// Create the token
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	// Set some claims
	token.Claims["iss"] = "https://accounts.google.com"
	token.Claims["aud"] = "api_client_id"
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token.Claims["sub"] = "12345"
	token.Claims["name"] = "Joe Bloggs"
	token.Claims["email"] = "joebloggs@email.com"
	token.Claims["picture"] = "picture.jpg"
	token.Header["kid"] = "test_key"

	// Sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString(signKey)
	f, _ := ioutil.ReadFile(pubKeyPath)
	var certs map[string]string
	json.Unmarshal(f, &certs)
	a := NewGoogleAuth(certs, tokenIssuer, "api_client_id")

	// retrieve the token
	t, _ := a.Verify(tokenString)
	p, _ := PlayerFromToken(t)
	fmt.Println(p.Name)
	fmt.Println(p.Email)
	fmt.Println(p.GoogleId)
	fmt.Println(p.Picture)
	fmt.Println(t.Valid)

	// Output:
	// Joe Bloggs
	// joebloggs@email.com
	// 12345
	// picture.jpg
	// true
}
