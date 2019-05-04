package keys

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

//GetPrivateKey returns a rsa.PrivateKey struct from a file path
func GetPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

//GetPublicKey returns a rsa.PublicKey struct from a file path
func GetPublicKey(filepath string) (*rsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
