package token_test

import (
	"io/ioutil"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/51st-state/api/pkg/token"
	"github.com/51st-state/api/test"
)

func TestNew(t *testing.T) {
	tok := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
	}, "1234", 0)

	if tok.Data().UUID != "1234" {
		t.Fatal("invalid uuid casted")
	}

	if tok.Data().Type != 0 {
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
