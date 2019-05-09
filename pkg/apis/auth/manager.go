package auth

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/51st-state/api/pkg/apis/user"
	"github.com/51st-state/api/pkg/token"
	jwt "github.com/dgrijalva/jwt-go"
)

// Manager for authenticating a user
type Manager struct {
	pK   *rsa.PrivateKey
	user *user.GRPCClient
}

// NewManager for user authentication
func NewManager(prvKey *rsa.PrivateKey, u *user.GRPCClient) *Manager {
	return &Manager{
		prvKey,
		u,
	}
}

var (
	errInvalidEmailFormat = errors.New("invalid email format")
	emailRegexp           = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type incompletePassword struct {
	password string
}

func (i *incompletePassword) Password() string {
	return i.password
}

// LoginWCFUser logs a user in with their connected wcf user credentials
func (m *Manager) LoginWCFUser(ctx context.Context, c Credentials) (*Token, error) {
	var (
		err  error
		info *user.WCFUserInfo
	)

	if emailRegexp.MatchString(c.Username()) {
		info, err = m.user.GetWCFInfoByEmail(ctx, c.Username())
		if err != nil {
			return nil, err
		}
	} else {
		info, err = m.user.GetWCFInfoByUsername(ctx, c.Username())
		if err != nil {
			return nil, err
		}
	}

	u, err := m.user.GetByWCFUserID(ctx, info.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			u, err = m.user.Create(ctx, user.NewIncomplete(info.UserID, "", false))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := m.user.CheckPassword(ctx, u, c); err != nil {
		return nil, err
	}

	aT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).UnixNano(),
		Audience:  "default",
	}, &token.User{
		ID:   u.UUID(),
		Type: "user",
	})

	rT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).UnixNano(),
		Audience:  "auth/refresh",
	}, &token.User{
		ID:   u.UUID(),
		Type: "user",
	})

	return &Token{
		m.pK,
		aT,
		rT,
	}, nil
}

// LoginOnGameServer logs in a user whereas the server fetches the client token
func (m *Manager) LoginOnGameServer(ctx context.Context, c ServerCredentials) (*Token, error) {
	return nil, nil
}

// RefreshToken returns a new access and refresh token
func (m *Manager) RefreshToken(ctx context.Context, accessToken token.Token, refreshToken token.Token) (*Token, error) {
	return nil, nil
}
