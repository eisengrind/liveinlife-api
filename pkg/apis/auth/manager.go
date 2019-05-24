package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/51st-state/api/pkg/problems"

	"github.com/51st-state/api/pkg/apis/serviceaccount/key"
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
	saKey             key.Manager
}

// NewManager for user authentication
func NewManager(prvKey *rsa.PrivateKey, r Repository, u user.Manager, v *recaptcha.Verifier, saKey key.Manager) *Manager {
	return &Manager{
		prvKey,
		r,
		u,
		v,
		saKey,
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
// A recaptcha login is only available for default user logins
func (m *Manager) RecaptchaLogin(ctx context.Context, c Credentials) (*Token, error) {
	responseToken := recaptchaRespFromCtx(ctx)
	verifyResp, err := m.recaptchaVerifier.Verify(responseToken, "")
	if err != nil {
		return nil, err
	}

	/*if !verifyResp.Success {
		return nil, errors.New("recaptcha could not be verified")
	}*/

	return m.loginUser(ctx, c)
}

var errTooManyAttempts = problems.New("too many login attempts", "a login request has to provide a recaptcha", 425)

const serviceAccountLoginName = "_json_key"

// Login logs a user in with their connected wcf user credentials
func (m *Manager) Login(ctx context.Context, c Credentials) (*Token, error) {
	if c.Name() == serviceAccountLoginName {
		return m.loginServiceAccount(ctx, c)
	}

	return m.loginUser(ctx, c)
}

func (m *Manager) loginUser(ctx context.Context, c Credentials) (*Token, error) {
	info, err := m.user.GetWCFInfo(ctx, c.Name())
	if err != nil {
		return nil, err
	}

	u, err := m.user.GetByWCFUserID(ctx, info.UserID)
	if err == user.ErrNotFound {
		u, err = m.user.Create(ctx, user.NewIncomplete(info.UserID, "", "", "", false))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	attempts, err := m.repo.LoginAttemptsCountSince(
		ctx,
		fmt.Sprintf("user/%s", u.UUID()),
		time.Now().Add(-(time.Hour * 24)),
	)
	if err != nil {
		return nil, err
	}

	if attempts > 0 {
		return nil, errTooManyAttempts
	}

	if err := m.user.CheckPassword(ctx, u, c); err != nil {
		// TODO: get exact status code of the error since the error also could be a timeout error
		if err := m.repo.AddLoginAttempt(
			ctx,
			fmt.Sprintf("user/%s", u.UUID()),
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

var errServiceAccountGUIDNotEqual = errors.New("the service account guid's are not equal")

func (m *Manager) loginServiceAccount(ctx context.Context, c Credentials) (*Token, error) {
	var jsonKey key.ClientKey
	if err := json.Unmarshal([]byte(c.Password()), &jsonKey); err != nil {
		return nil, err
	}

	k, err := m.saKey.Get(ctx, key.NewIdentifier(jsonKey.GUID))
	if err != nil {
		return nil, err
	}

	if k.Data().ServiceAccountGUID != jsonKey.ServiceAccountGUID {
		return nil, errServiceAccountGUIDNotEqual
	}

	if err := validateKeypair(jsonKey.GUID, jsonKey.PrivateKey, k.Data().PublicKey); err != nil {
		return nil, err
	}

	aT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Audience:  "default",
	}, &token.User{
		ID:   jsonKey.ServiceAccountGUID,
		Type: "service_account",
	})

	rT := token.New(&jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		Audience:  "auth/refresh",
	}, &token.User{
		ID:   jsonKey.ServiceAccountGUID,
		Type: "service_account",
	})

	return &Token{
		m.pK,
		aT,
		rT,
	}, nil
}

var errKeypairNotMatching = errors.New("the given keypair does not match")

func validateKeypair(verificationText string, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	h := sha256.New()
	h.Write([]byte(verificationText))
	digest := h.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, sig)
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
