package key

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/encode"
	"github.com/51st-state/api/pkg/rbac"
	rbacMiddleware "github.com/51st-state/api/pkg/rbac/middleware"
	"github.com/51st-state/api/pkg/token"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// MakeGetEndpoint creates a new http endpoint for the return of service account keys
func MakeGetEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		guid := chi.URLParam(r, "guid")
		return m.Get(ctx, &identifier{guid})
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.keys.get"))).
		HandlerFunc(l)
}

// MakeCreateEndpoint creates a new http endpoint for the creation of service account keys
func MakeCreateEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		serviceAccountGUID := chi.URLParam(r, "guid")

		inc := NewIncomplete("", "")
		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}
		inc.Data().ServiceAccountGUID = serviceAccountGUID

		return m.Create(ctx, inc)
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.keys.create"))).
		HandlerFunc(l)
}

// MakeSetEndpoint creates a new http endpoint for updating the service account key info
func MakeSetEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		guid := chi.URLParam(r, "guid")

		inc := NewIncomplete("", "")
		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}
		inc.Data().ServiceAccountGUID = ""
		inc.Data().PublicKey = nil

		return struct{}{}, m.Update(ctx, &complete{
			&identifier{guid},
			inc,
		})
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.keys.set"))).
		HandlerFunc(l)
}

// MakeDeleteEndpoint create a new http endpoint to delete service account keys
func MakeDeleteEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		guid := chi.URLParam(r, "guid")
		return struct{}{}, m.Delete(ctx, &identifier{guid})
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.keys.delete"))).
		HandlerFunc(l)
}
