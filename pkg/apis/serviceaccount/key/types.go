package key

//go:generate counterfeiter -o ./mocks/identifier.go . Identifier

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// Identifier of a service account keypair
type Identifier interface {
	GUID() string
}

type identifier struct {
	guid string
}

// NewIdentifier creates a new identifier object
func NewIdentifier(g string) Identifier {
	return &identifier{g}
}

func (i *identifier) GUID() string {
	return i.guid
}

// Provider provides data of a service account keypair
type Provider interface {
	Data() *data
}

// Incomplete represents an incomplete service account keypair
type Incomplete interface {
	Provider
}

// Complete represents a complete service account keypair object
type Complete interface {
	Identifier
	Incomplete
}

type complete struct {
	Identifier
	Incomplete
}

type data struct {
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	PublicKey          *rsa.PublicKey `json:"-"`
	ServiceAccountGUID string         `json:"service_account_guid"`
}

// NewIncomplete creates a new incomplete service account keypair object
func NewIncomplete(name, description string) Incomplete {
	return &data{
		name,
		description,
		nil,
		"",
	}
}

func (d *data) Data() *data {
	return d
}

func (d *data) SetName(to string) *data {
	d.Name = to
	return d
}

func (d *data) SetDescription(to string) *data {
	d.Description = to
	return d
}

// ClientKey represents the private key a client receives when creating a service account keypair
type ClientKey struct {
	ServiceAccountGUID string          `json:"service_account_guid"`
	GUID               string          `json:"guid"`
	PrivateKey         *rsa.PrivateKey `json:"-"`
}

// MarshalJSON marshals the client key to presentable data
func (k *ClientKey) MarshalJSON() ([]byte, error) {
	var buf = bytes.NewBuffer([]byte{})

	if err := pem.Encode(buf, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k.PrivateKey),
	}); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		ServiceAccountGUID string `json:"service_account_guid"`
		GUID               string `json:"guid"`
		PrivateKey         string `json:"private_key"`
	}{
		k.ServiceAccountGUID,
		k.GUID,
		string(b),
	})
}

var (
	errInvalidPEMKey       = errors.New("invalid pem key")
	errInvalidPEMBlockType = errors.New("invalid pem block type")
	errPrivateKeyInvalid   = errors.New("private key is invalid")
)

// UnmarshalJSON unmarshals incoming JSON
func (k *ClientKey) UnmarshalJSON(b []byte) error {
	var key struct {
		ServiceAccountGUID string `json:"service_account_guid"`
		GUID               string `json:"guid"`
		PrivateKey         string `json:"private_key"`
	}

	if err := json.Unmarshal(b, &key); err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(key.PrivateKey))
	if block == nil {
		return errInvalidPEMKey
	}

	if block.Type != "PRIVATE KEY" {
		return errInvalidPEMBlockType
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	// this statement is probably unnecessary due to ParsePKCS1Privatekey([...])
	// but for security reasons there is a double check
	// TODO: check if this statement can be removed
	if privateKey.Validate() != nil {
		return errPrivateKeyInvalid
	}

	k.PrivateKey = privateKey
	k.GUID = key.GUID
	k.ServiceAccountGUID = key.ServiceAccountGUID

	return nil
}
