package role

import (
	"context"
	"crypto/rsa"
	"net/http"
	"regexp"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/rbac"
	"github.com/51st-state/api/pkg/token"

	"github.com/51st-state/api/pkg/encode"
	"go.uber.org/zap"
)

// MakeGetEndpoint for the role service
func MakeGetEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeSetEndpoint for the role service
func MakeSetEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeCreateEndpoint for the role service
func MakeCreateEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

var systemRoleRegexp = regexp.MustCompile(`^(system\/)?[a-z0-9-_]+$`)

// MakeDeleteEndpoint for the role service
func MakeDeleteEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}
