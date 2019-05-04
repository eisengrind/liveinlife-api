package keys_test

import (
	"testing"

	"github.com/51st-state/api/pkg/keys"
)

func TestGetPrivateKey(t *testing.T) {
	if _, err := keys.GetPrivateKey("../test"); err == nil {
		t.Fatal("i know that this is the wrong path...")
	}

	if _, err := keys.GetPrivateKey("../../test/testPublicKey.pem"); err == nil {
		t.Fatal("Though, this is no real key file")
	}

	if _, err := keys.GetPrivateKey("../../test/testPrivateKey.pem"); err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetPublicKey(t *testing.T) {
	if _, err := keys.GetPublicKey("../test"); err == nil {
		t.Fatal("i know that this is the wrong path...")
	}

	if _, err := keys.GetPublicKey("../../test/testPrivateKey.pem"); err == nil {
		t.Fatal("Though, this is no real key file")
	}

	if _, err := keys.GetPublicKey("../../test/testPublicKey.pem"); err != nil {
		t.Fatal(err.Error())
	}
}
