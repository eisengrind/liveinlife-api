package auth

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/51st-state/api/pkg/problems"

	"github.com/51st-state/api/pkg/apis/user"
	"github.com/51st-state/api/pkg/recaptcha"
	"github.com/51st-state/api/pkg/token"
	jwt "github.com/dgrijalva/jwt-go"
)

// Manager for authenticating a user
type Manager struct {
	pK                *rsa.PrivateKey
	repo              Repository
	user              user.Manager
	recaptchaVerifier *recaptcha.Verifier
}

// NewManager for user authentication
func NewManager(prvKey *rsa.PrivateKey, r Repository, u user.Manager, v *recaptcha.Verifier) *Manager {
	return &Manager{
		prvKey,
		r,
		u,
		v,
	}
}

var (
	errInvalidEmailFormat = errors.New("invalid email format")
)

type incompletePassword struct {
	password string
}

func (i *incompletePassword) Password() string {
	return i.password
}

// RecaptchaLogin logs a user in with a check for recaptcha
func (m *Manager) RecaptchaLogin(ctx context.Context, c Credentials) (*Token, error) {
	responseToken := recaptchaRespFromCtx(ctx)
	verifyResp, err := m.recaptchaVerifier.Verify(responseToken, "")
	if err != nil {
		return nil, err
	}

	if len(verifyResp.ErrorCodes) != 0 {
		return nil, errors.New("recaptcha validation errored")
	}

	return m.login(ctx, c)
}

var errTooManyAttempts = problems.New("too many login attempts", "a login request has to provide a recaptcha", 425)

// Login logs a user in with their connected wcf user credentials
func (m *Manager) Login(ctx context.Context, c Credentials) (*Token, error) {
	attempts, err := m.repo.LoginAttemptsCountSince(
		ctx,
		fmt.Sprintf("user/%s", c.Username()),
		time.Now().Add(-(time.Hour * 24)),
	)
	if err != nil {
		return nil, err
	}

	if attempts > 0 {
		return nil, errTooManyAttempts
	}

	return m.login(ctx, c)
}

func (m *Manager) login(ctx context.Context, c Credentials) (*Token, error) {
	info, err := m.user.GetWCFInfo(ctx, c.Username())
	if err != nil {
		return nil, err
	}

	u, err := m.user.GetByWCFUserID(ctx, info.UserID)
	if err == sql.ErrNoRows {
		u, err = m.user.Create(ctx, user.NewIncomplete(info.UserID, "", false))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if err := m.user.CheckPassword(ctx, u, c); err != nil {
		// TODO: get exact status code of the error since the error also could be a timeout error
		if err := m.repo.AddLoginAttempt(
			ctx,
			fmt.Sprintf("user/%s", c.Username()),
			time.Now(),
		); err != nil {
			return nil, err
		}

		return nil, err
	}

	aT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Audience:  "default",
	}, &token.User{
		ID:   u.UUID(),
		Type: "user",
	})

	rT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
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

// RefreshToken returns a new access and refresh token
func (m *Manager) RefreshToken(ctx context.Context, accessToken token.Token, refreshToken token.Token) (*Token, error) {
	if accessToken.Data().Audience != "default" {
		return nil, errors.New("access token has an invalid audience")
	}

	if refreshToken.Data().Audience != "auth/refresh" {
		return nil, errors.New("refresh token has an invalid audience")
	}

	if accessToken.Data().User.String() != refreshToken.Data().User.String() {
		return nil, errors.New("the UUIDs are not equal")
	}

	aT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).UnixNano(),
		Audience:  "default",
	}, &token.User{
		ID:   accessToken.Data().User.ID,
		Type: accessToken.Data().User.Type,
	})

	rT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).UnixNano(),
		Audience:  "auth/refresh",
	}, &token.User{
		ID:   accessToken.Data().User.ID,
		Type: accessToken.Data().User.Type,
	})

	return &Token{
		m.pK,
		aT,
		rT,
	}, nil
}
