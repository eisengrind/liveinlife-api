package auth

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/51st-state/api/pkg/token"
)

// Credentials to login and to fetch access and refresh token
type Credentials interface {
	Username() string
	Password() string
}

type credentials struct {
	username string
	password string
}

func newCredentials(u, p string) Credentials {
	return &credentials{u, p}
}

func (c *credentials) Username() string {
	return c.username
}

func (c *credentials) Password() string {
	return c.password
}

func (c *credentials) UnmarshalJSON(b []byte) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(b, &req); err != nil {
		return err
	}

	c.username = req.Username
	c.password = req.Password

	return nil
}

// ServerCredentials for logins on the server
type ServerCredentials interface {
	GameSerialHash() string
	Password() string
}

type serverCredentials struct {
	hash     string
	password string
}

func newServerCredentials(h, p string) ServerCredentials {
	return &serverCredentials{h, p}
}

func (c *serverCredentials) GameSerialHash() string {
	return c.hash
}

func (c *serverCredentials) Password() string {
	return c.password
}

// Token to return to the client
type Token struct {
	pK           *rsa.PrivateKey
	accessToken  token.Token
	refreshToken token.Token
}

// MarshalJSON for a token
func (t *Token) MarshalJSON() ([]byte, error) {
	aT, err := t.accessToken.String(t.pK)
	if err != nil {
		return nil, err
	}

	rT, err := t.refreshToken.String(t.pK)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		aT,
		rT,
	})
}
