package token

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// Token authenticates a consumer of the api
type Token interface {
	// Token returns the initial JWT-Token
	Token() *jwt.Token
	// String returns a signed and readable JWT-Token
	String(key *rsa.PrivateKey) (string, error)
	// Data of the token
	Data() *data
}

type tokenInfo struct {
	ID   string `json:"id"`
	Type uint8  `json:"type"`
}

type data struct {
	*tokenInfo
	*jwt.StandardClaims
}

type token struct {
	token *jwt.Token
}

func (t *token) Token() *jwt.Token {
	return t.token
}

func (t *token) String(key *rsa.PrivateKey) (string, error) {
	return t.token.SignedString(key)
}

func (t *token) Data() *data {
	return t.token.Claims.(*data)
}

// New creates a new api-friendly api token.
// Available types are 0 = user-token, 1 = user-refresh-token.
// More types can follow - e.g. specific API-Token in addition to user token.
func New(c *jwt.StandardClaims, id string, t uint8) Token {
	return &token{
		token: jwt.NewWithClaims(jwt.SigningMethodRS512, &data{
			&tokenInfo{
				ID:   id,
				Type: t,
			},
			c,
		}),
	}
}

// NewFromString parses a token object from a jwt token string
func NewFromString(pK rsa.PublicKey, t string) (Token, error) {
	tok, err := jwt.ParseWithClaims(
		t,
		&data{},
		func(_ *jwt.Token) (interface{}, error) {
			return pK, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return Convert(tok)
}

// Convert a standard JWT-Token implementing necessary
// data to a Token interface.
func Convert(t *jwt.Token) (Token, error) {
	if _, ok := t.Claims.(*data); !ok {
		return nil, errors.New("token does not implement data")
	}

	return &token{
		token: t,
	}, nil
}
