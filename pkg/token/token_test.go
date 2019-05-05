package token_test

import (
	"testing"
	"time"

	"github.com/51st-state/api/pkg/keys"

	"github.com/51st-state/api/pkg/token"
	"github.com/51st-state/api/test"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestNew(t *testing.T) {
	tok := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
	}, &token.User{
		ID:   "username",
		Type: "user",
	})

	if tok.Data().User.ID != "username" {
		t.Fatal("invalid id casted")
	}

	if tok.Data().User.Type != "user" {
		t.Fatal("invalid type casted")
	}

	privateKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := tok.String(privateKey); err != nil {
		t.Fatal(err.Error())
	}
}

func TestNewFromString(t *testing.T) {
	tok := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
	}, &token.User{
		ID:   "username",
		Type: "user",
	})

	privateKey, err := keys.GetPrivateKey(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}

	tokStr, err := tok.String(privateKey)
	if err != nil {
		t.Fatal(err.Error())
	}

	publicKey, err := keys.GetPublicKey(test.GetTestPublicKey())
	if err != nil {
		t.Fatal(err.Error())
	}

	tok, err = token.NewFromString(publicKey, tokStr)
	if err != nil {
		t.Fatal(err.Error())
	}
}
