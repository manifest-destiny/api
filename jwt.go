package api

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/square/go-jose"
)

var (
	googleTokenIssuers = []string{"accounts.google.com", "https://accounts.google.com"}
	googleJWKSetURL    = "https://www.googleapis.com/oauth2/v3/certs"
)

// Timestamp custom time.Time type with MarshalJSON and UnmarshalJSON methods.
type Timestamp struct {
	time.Time
}

// MarshalJSON handles Timestamp to byte array conversion.
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	fs := fmt.Sprint(ts)

	return []byte(fs), nil
}

// UnmarshalJSON handles JSON timestamp string to Timestamp type.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	t.Time = time.Unix(int64(ts), 0)

	return nil
}

// Validator defines Claims keys to validate.
type Validator struct {
	Issuers        []string
	Subjects       []string
	Audiences      []string
	IDs            []string
	CheckExpired   bool
	CheckNotBefore bool
}

// Claimer represents a set of JWT claims.
type Claimer interface {
	Valid() error
}

// Claims represents a JWT claims set. Includes standard keys with
// additional Google claims.
type Claims struct {
	Validator      *Validator
	Issuer         string    `json:"iss,omitempty"`
	Subject        string    `json:"sub,omitempty"`
	Audience       string    `json:"aud,omitempty"`
	ExpirationTime Timestamp `json:"exp,omitempty"`
	NotBefore      Timestamp `json:"nbf,omitempty"`
	IssuedAt       Timestamp `json:"iat,omitempty"`
	ID             string    `json:"jti,omitempty"`
}

// GoogleIdentityClaims represents Google JWT claims and extends Claims.
type GoogleIdentityClaims struct {
	*Claims
	Email               string `json:"email,omitempty"`
	EmailVerified       bool   `json:"email_verified,omitempty"`
	Name                string `json:"name,omitempty"`
	GivenName           string `json:"given_name,omitempty"`
	FamilyName          string `json:"family_mame,omitempty"`
	Locale              string `json:"locale,omitempty"`
	Picture             string `json:"picture,omitempty"`
	AuthorizedPresenter string `json:"azp,omitempty"`
	AccessTokenHash     string `json:"at_hash,omitempty"`
}

// Valid validates a Claims set.
func (c *Claims) Valid() error {
	validClaim := func(name, claim string, expected []string) (bool, error) {
		if len(expected) == 0 {
			return true, nil
		}
		for _, expect := range expected {
			if expect == claim {
				return true, nil
			}
		}

		return false, fmt.Errorf("Invalid \"%s\" claim: \"%s\"", name, claim)
	}

	if ok, err := validClaim("iss", c.Issuer, c.Validator.Issuers); !ok {
		return err
	}
	if ok, err := validClaim("sub", c.Subject, c.Validator.Subjects); !ok {
		return err
	}
	if ok, err := validClaim("aud", c.Audience, c.Validator.Audiences); !ok {
		return err
	}
	if ok, err := validClaim("jti", c.ID, c.Validator.IDs); !ok {
		return err
	}

	now := time.Now()
	if c.Validator.CheckExpired && !now.Before(c.ExpirationTime.Time) {
		return fmt.Errorf("Claim expired")
	}
	if c.Validator.CheckNotBefore && !now.After(c.NotBefore.Time) {
		return fmt.Errorf("Claim not yet valid")
	}

	return nil
}

// JWSVerifyer struct verifies JSON Web Signatures against a JWK Set.
type JWSVerifyer struct {
	*jose.JsonWebKeySet
	keySetURL string
	expires   time.Time
}

// NewJWSVerifyer constructs a JWSVerifyer.
func NewJWSVerifyer(url string) (*JWSVerifyer, error) {
	v := &JWSVerifyer{keySetURL: url}
	err := v.setKeySet()
	if err != nil {
		return v, err
	}

	return v, nil
}

func (v *JWSVerifyer) setKeySet() error {

	resp, err := http.Get(v.keySetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jwkSet := &jose.JsonWebKeySet{}
	err = json.Unmarshal(body, &jwkSet)
	if err != nil {
		return err
	}

	header := resp.Header.Get("Expires")
	expires, err := time.Parse(time.RFC1123, header)
	if err != nil {
		return err
	}

	v.JsonWebKeySet = jwkSet
	v.expires = expires

	return nil
}

// UpdateKeySet updates the JWK Set.
func (v *JWSVerifyer) UpdateKeySet() error {
	now := time.Now()
	if now.After(v.expires) {
		err := v.setKeySet()
		if err != nil {
			return err
		}
	}

	return nil
}

// Verify verifies token against signing key and returns JWT claims.
func (v *JWSVerifyer) Verify(token string, c Claimer) error {

	err := v.UpdateKeySet()
	if err != nil {
		return err
	}

	obj, err := jose.ParseSigned(token)
	if err != nil {
		return err
	}

	if len(obj.Signatures) < 1 {
		return fmt.Errorf("need a token signature, found %d", len(obj.Signatures))
	}

	kid := obj.Signatures[0].Header.KeyID
	keys := v.Key(kid)
	if len(keys) == 0 {
		return fmt.Errorf("need at least one key with ID: %s", kid)
	}

	var key *rsa.PublicKey
	switch keys[0].Key.(type) {
	case *rsa.PublicKey:
		key = keys[0].Key.(*rsa.PublicKey)
	default:
		return fmt.Errorf("unsupported key type")
	}
	output, err := obj.Verify(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(output, c); err != nil {
		return err
	}

	return nil
}

// TokenValidator convenience container for JWSVerifyer and Claimer.
type TokenValidator struct {
	Claims     *Validator
	Signatures *JWSVerifyer
}

// GoogleTokenValidator constructor for a Google TokenValidator.
func GoogleTokenValidator(aud ...string) (*TokenValidator, error) {
	kv, err := NewJWSVerifyer(googleJWKSetURL)
	if err != nil {
		return &TokenValidator{}, err
	}

	cv := &Validator{
		Issuers:      googleTokenIssuers,
		Audiences:    aud,
		CheckExpired: true,
	}

	return &TokenValidator{
		cv,
		kv,
	}, nil
}
