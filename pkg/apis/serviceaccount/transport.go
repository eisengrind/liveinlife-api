package serviceaccount

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/token"

	"github.com/51st-state/api/pkg/rbac"
	rbacMiddleware "github.com/51st-state/api/pkg/rbac/middleware"

	"github.com/51st-state/api/pkg/encode"
	"go.uber.org/zap"
)

// MakeGetEndpoint creates a new http endpoint to return a service account
func MakeGetEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		guid := chi.URLParam(r, "guid")
		return m.Get(ctx, &identifier{guid})
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.get"))).
		HandlerFunc(l)
}

// MakeUpdateEndpoint creates a new http endpoint updates an already created service account
func MakeUpdateEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		guid := chi.URLParam(r, "guid")

		inc := NewIncomplete("", "")

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return struct{}{}, m.Update(ctx, &complete{
			&identifier{guid},
			inc,
		})
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.set"))).
		HandlerFunc(l)
}

// MakeCreateEndpoint creates a new http endpoint to create a service account
func MakeCreateEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		inc := NewIncomplete("", "")

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return m.Create(ctx, inc)
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.create"))).
		HandlerFunc(l)
}

// MakeDeleteEndpoint creates a new http endpoint to delete a service account
func MakeDeleteEndpoint(l *zap.Logger, e encode.Encoder, publicKey *rsa.PublicKey, m Manager, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		guid := chi.URLParam(r, "guid")
		return struct{}{}, m.Delete(ctx, &identifier{guid})
	}).
		WithBefore(token.NewMiddleware(*publicKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("serviceaccounts.delete"))).
		HandlerFunc(l)
}
