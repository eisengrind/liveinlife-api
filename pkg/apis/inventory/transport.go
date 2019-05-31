package inventory

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

// MakeGetEndpoint creates a http endpoint to retrieve an inventory object
func MakeGetEndpoint(l *zap.Logger, m Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		return m.Get(ctx, &identifier{chi.URLParam(r, "guid")})
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("inventory.get"))).
		HandlerFunc(l)
}

// MakeCreateEndpoint creates a http endpoint to create an inventory object
func MakeCreateEndpoint(l *zap.Logger, m Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		inc := NewIncomplete(make([]*Item, 0))

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return m.Create(ctx, inc)
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("inventory.create"))).
		HandlerFunc(l)
}

// MakeAddItemEndpoint creates a http endpoint to add an item to an inventory object
func MakeAddItemEndpoint(l *zap.Logger, m Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		id := &identifier{chi.URLParam(r, "guid")}

		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			return nil, err
		}

		return struct{}{}, m.AddItem(ctx, id, &item)
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("inventory.item.add"))).
		HandlerFunc(l)
}

// MakeRemoveItemEndpoint creates a http endpoint to remove an item from an inventory object
func MakeRemoveItemEndpoint(l *zap.Logger, m Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		id := &identifier{chi.URLParam(r, "guid")}

		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			return nil, err
		}

		return struct{}{}, m.RemoveItem(ctx, id, &item)
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("inventory.item.remove"))).
		HandlerFunc(l)
}

// MakeDeleteEndpoint creates a http endpoint to delete an inventory object
func MakeDeleteEndpoint(l *zap.Logger, m Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		return struct{}{}, m.Delete(ctx, &identifier{chi.URLParam(r, "guid")})
	}).
		WithBefore(token.NewMiddleware(pubKey)).
		WithBefore(rbacMiddleware.NewRulecheck(rb, rbac.Rule("inventory.delete"))).
		HandlerFunc(l)
}
