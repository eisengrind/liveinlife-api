package role

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/rbac"
	rbacMiddleware "github.com/51st-state/api/pkg/rbac/middleware"
	"github.com/51st-state/api/pkg/token"

	"github.com/51st-state/api/pkg/encode"
	"go.uber.org/zap"
)

// MakeGetEndpoint for the role service
func MakeGetEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		id := chi.URLParam(r, "id")
		return m.Get(ctx, newIdentifier(rbac.RoleID(id)))
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("roles.get"))).
		HandlerFunc(l)
}

// MakeSetEndpoint for the role service
func MakeSetEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		id := chi.URLParam(r, "id")
		inc := NewIncomplete("", "", make(rbac.RoleRules, 0))

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return struct{}{}, m.Set(
			ctx,
			&complete{
				newIdentifier(rbac.RoleID(id)),
				inc,
			},
		)
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("roles.set"))).
		HandlerFunc(l)
}

// MakeCreateEndpoint for the role service
func MakeCreateEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		id := chi.URLParam(r, "id")
		inc := NewIncomplete("", "", rbac.RoleRules{})

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return struct{}{}, m.Create(ctx, &complete{
			newIdentifier(rbac.RoleID(id)),
			inc,
		})
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("roles.create"))).
		HandlerFunc(l)
}

var systemRoleRegexp = regexp.MustCompile(`^(system\/)?[a-z0-9-_]+$`)

// MakeDeleteEndpoint for the role service
func MakeDeleteEndpoint(l *zap.Logger, m Manager, e encode.Encoder, pubKey rsa.PublicKey, rb rbac.Control) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		id := chi.URLParam(r, "id")
		return struct{}{}, m.Delete(ctx, newIdentifier(rbac.RoleID(id)))
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("roles.delete"))).
		HandlerFunc(l)
}
