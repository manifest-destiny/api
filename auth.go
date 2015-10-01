package main

import (
	"crypto/rsa"
	"errors"
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/manifest-destiny/api/player"
)

const (
	tokenIssuer      = "accounts.google.com"
	googlePubKeyPath = "keys/google.pub"
	webAppClientId   = "801574721267-8ocanqgcgln83r5s2bdpk5imu78r2ouk.apps.googleusercontent.com"
)

type GoogleAuth struct {
	verifyKeys map[string]*rsa.PublicKey
	clientIds  []string
	issuer     string
}

type GoogleToken struct {
	*jwt.Token
	Valid bool
}

func NewGoogleAuth(certs map[string]string, issuer string, clientIds ...string) *GoogleAuth {

	verifyKeys := make(map[string]*rsa.PublicKey)

	for key, cert := range certs {
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			log.Fatal(err)
		}
		verifyKeys[key] = verifyKey
	}

	return &GoogleAuth{
		verifyKeys: verifyKeys,
		clientIds:  clientIds,
		issuer:     issuer,
	}
}

func (a *GoogleAuth) Auth(tokenString string) (*player.Player, error) {
	t, err := a.Verify(tokenString)
	if err != nil {
		return &player.Player{}, err
	}
	return PlayerFromToken(t)
}

func (a *GoogleAuth) Verify(tokenString string) (*GoogleToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); !ok {
			return nil, errors.New("Cannot get kid from header")
		} else {
			if pubKey, ok := a.verifyKeys[kid]; ok {
				return pubKey, nil
			} else {
				return &rsa.PublicKey{}, errors.New("Uknown key ID")
			}
		}
	})

	gToken := &GoogleToken{
		token,
		token.Valid,
	}

	if err != nil || !gToken.Valid {
		return gToken, err
	}

	err = a.validateClaims(gToken)
	if err != nil {
		gToken.Valid = false
	}
	return gToken, err
}

func (a *GoogleAuth) validateClaims(token *GoogleToken) error {
	err := a.validateIssuer(token)
	if err != nil {
		return err
	}
	return a.validateClient(token)
}

func (a *GoogleAuth) validateIssuer(token *GoogleToken) error {
	// Strip protocol scheme from issuer (if it exists) https://accounts.google.com
	if str, ok := token.Claims["iss"].(string); ok {
		iss := strings.TrimLeft(str, "https://")
		if iss != a.issuer {
			return errors.New("Invalid issuer")
		}
	} else {
		return errors.New("Cannot get iss from claims")
	}

	return nil
}

func (a *GoogleAuth) validateClient(token *GoogleToken) error {

	hasClient := func(aud string) bool {
		for _, clientId := range a.clientIds {
			if clientId == aud {
				return true
			}
		}
		return false
	}

	// check if client id is valid
	if str, ok := token.Claims["aud"].(string); ok {
		if hasClient(str) == false {
			return errors.New("Invalid API client")
		}
	} else {
		return errors.New("Cannot get aud from claims")
	}

	return nil
}

func PlayerFromToken(token *GoogleToken) (*player.Player, error) {
	p := &player.Player{}
	if sub, ok := token.Claims["sub"].(string); ok {
		p.GoogleId = sub
	} else {
		return p, errors.New("Cannot parse google ID")
	}
	if email, ok := token.Claims["email"].(string); ok {
		p.Email = email
	}
	if name, ok := token.Claims["name"].(string); ok {
		p.Name = name
	}
	if pic, ok := token.Claims["picture"].(string); ok {
		p.Picture = pic
	}
	return p, nil
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
