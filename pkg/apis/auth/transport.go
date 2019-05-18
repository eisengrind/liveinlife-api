package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/encode"
	"github.com/51st-state/api/pkg/token"
	"go.uber.org/zap"
)

// MakeRecaptchaLoginEndpoint creates a new http endpoint for a recaptcha login
func MakeRecaptchaLoginEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		creds := newCredentials("", "")
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			return nil, err
		}

		return m.RecaptchaLogin(ctx, creds)
	}).
		WithBefore(populateRecaptchaResponseToken()).
		HandlerFunc(l)
}

// MakeLoginEndpoint creates a new http endpoint for a normal login
func MakeLoginEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		creds := newCredentials("", "")
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			return nil, err
		}

		return m.Login(ctx, creds)
	}).
		HandlerFunc(l)
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// MakeRefreshTokenEndpoint creates a new http endpoint for refreshing
// an access token with a given refresh token
func MakeRefreshTokenEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		var req refreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}

		accessToken, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		refreshToken, err := token.NewFromString(&pubKey, req.RefreshToken)
		if err != nil {
			return nil, err
		}

		return m.RefreshToken(ctx, accessToken, refreshToken)
    }).
        WithBefore(func(ctx context.Context, r *http.Request) (context.Context, error) {
            l.Info("access token: ", zap.String("access_token", r.Header.Get("Authorization")))
            return ctx, nil
        }).
		WithBefore(token.NewMiddleware(pubKey)).
		HandlerFunc(l)
}
