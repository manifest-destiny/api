package api

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/square/go-jose"
)

var now = time.Now()

var validationCases = []struct {
	*Claims
	Pass bool
	Name string
}{
	{
		&Claims{
			Validator:      &Validator{},
			Issuer:         "iss-val",
			Subject:        "sub-val",
			Audience:       "aud-val",
			ExpirationTime: Timestamp{now.Add(time.Hour * -1)},
			NotBefore:      Timestamp{now.Add(time.Hour * 1)},
		},
		true,
		"zeroed validator",
	},
	{
		&Claims{
			Validator:      &Validator{CheckExpired: true},
			ExpirationTime: Timestamp{now.Add(time.Hour * -1)},
		},
		false,
		"claim expired",
	},
	{
		&Claims{
			Validator: &Validator{CheckNotBefore: true},
			NotBefore: Timestamp{now.Add(time.Hour * 1)},
		},
		false,
		"claim not yet valid",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"issuer"}},
			Issuer:    "iss-val",
		},
		false,
		"invalid issuer",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"issuer", "iss-val"}},
			Issuer:    "iss-val",
		},
		true,
		"valid issuer",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"audience"}},
			Audience:  "aud-val",
		},
		false,
		"invalid audience",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"audience", "aud-val"}},
			Issuer:    "aud-val",
		},
		true,
		"valid audience",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"subject"}},
			Audience:  "sub-val",
		},
		false,
		"invalid subject",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"subject", "sub-val"}},
			Issuer:    "sub-val",
		},
		true,
		"valid subject",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"id"}},
			Audience:  "id-val",
		},
		false,
		"invalid id",
	},
	{
		&Claims{
			Validator: &Validator{Issuers: []string{"id", "id-val"}},
			Issuer:    "id-val",
		},
		true,
		"valid id",
	},
}

func TestClaimsValidate(t *testing.T) {
	for _, c := range validationCases {
		err := c.Claims.Valid()
		if err != nil && c.Pass {
			t.Errorf("Expected to pass validation for \"%s\", got error: %s", c.Name, err)
		} else if err == nil && !c.Pass {
			t.Errorf("Expected to fail validation for \"%s\", got a pass", c.Name)
		}
	}
}

func TestVerifyToken(t *testing.T) {
	// Create mock key set endpoint
	privateKey, _ := rsa.GenerateKey(rand.Reader, 512)
	jwkSet := jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{
			jose.JsonWebKey{
				Key:       privateKey.Public(),
				KeyID:     "1",
				Algorithm: "RSA256",
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(jwkSet)
		if err != nil {
			t.Error(err)
		}
		w.Header().Set("Expires", time.Now().Add(time.Hour*1).Format(time.RFC1123))
		w.Write(body)
	}))
	defer ts.Close()

	// Create JWT token with claims
	jwtClaims := `{"iss": "me"}`
	jwt, err := makeJwt(privateKey, jwtClaims, "1", "RSA256")
	if err != nil {
		t.Error(err)
	}
	// Create a JWS verifyer from JWK Set URL
	v, err := NewJWSVerifyer(ts.URL)
	if err != nil {
		t.Error(err)
	}
	c := &Claims{}
	// Verify JWT signature
	err = v.Verify(jwt, c)
	if err != nil {
		t.Error(err)
	}
	// Check round trip
	if c.Issuer != "me" {
		t.Error("expected issuser to be \"me\", got:", c.Issuer)
	}
}

func makeJwt(k *rsa.PrivateKey, claims, kid, alg string) (string, error) {
	jwk := jose.JsonWebKey{Key: k, KeyID: kid, Algorithm: alg}
	signer, err := jose.NewSigner(jose.RS256, &jwk)
	if err != nil {
		return "", err
	}
	obj, err := signer.Sign([]byte(claims))
	if err != nil {
		return "", err
	}

	return obj.CompactSerialize()
}
