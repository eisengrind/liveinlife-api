package test

import (
	"path/filepath"
	"runtime"
)

const (
	//PublicKeyPath is a public test sign file
	PublicKeyPath = "testPublicKey.pem"
	//PrivateKeyPath is private test sign file
	PrivateKeyPath = "testPrivateKey.pem"
)

func getPackageDir() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Dir(f)
}

//GetTestPublicKey returns the path to the test public key
func GetTestPublicKey() string {
	p, err := filepath.Abs(filepath.Join(getPackageDir(), PublicKeyPath))
	if err != nil {
		panic(err)
	}
	return p
}

//GetTestPrivateKey returns the path to the test private key
func GetTestPrivateKey() string {
	p, err := filepath.Abs(filepath.Join(getPackageDir(), PrivateKeyPath))
	if err != nil {
		panic(err)
	}
	return p
}
