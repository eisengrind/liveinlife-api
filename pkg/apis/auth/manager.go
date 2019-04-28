package auth

import (
	"context"
	"crypto/rsa"
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

// Login a user on their own
func (m *Manager) Login(ctx context.Context, c Credentials) (t *Token, err error) {
	var info *user.WCFUserInfo

	if emailRegexp.MatchString(c.Username()) {
		info, err = m.user.GetWCFInfoByEmail(ctx, c.Username())
		if err != nil {
			return
		}
	} else {
		info, err = m.user.GetWCFInfoByUsername(ctx, c.Username())
		if err != nil {
			return
		}
	}

	user, err := m.user.GetByWCFUserID(ctx, info.UserID)
	if err != nil {
		return
	}

	aT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).UnixNano(),
	}, user.UUID(), 0)

	rT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).UnixNano(),
	}, user.UUID(), 1)

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
