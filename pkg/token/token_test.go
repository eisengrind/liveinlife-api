package token_test

import (
	"io/ioutil"
	"testing"
	"time"

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

	b, err := ioutil.ReadFile(test.GetTestPrivateKey())
	if err != nil {
		t.Fatal(err.Error())
	}

	prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := tok.String(prvKey); err != nil {
		t.Fatal(err.Error())
	}
}
