package user

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/encode"
	"github.com/51st-state/api/pkg/rbac"
	"github.com/51st-state/api/pkg/token"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// MakeGetEndpoint for the user service
// API-Path: /users/{uuid}
func MakeGetEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.get")); err != nil {
			return nil, err
		}

		return m.Get(ctx, newIdentifier(chi.URLParam(r, "uuid")))
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeGetByGameSerialHashEndpoint for the user service
// API-Path: /users/hash/{hash}
func MakeGetByGameSerialHashEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.getByHash")); err != nil {
			return nil, err
		}

		return m.GetByGameSerialHash(ctx, chi.URLParam(r, "hash"))
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeCreateEndpoint for the user service
// API-Path: /users
func MakeCreateEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.create")); err != nil {
			return nil, err
		}

		inc := NewIncomplete(0, "", false)

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return m.Create(ctx, inc)
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeDeleteEndpoint for the user service
// API-Path: /users/{uuid}
func MakeDeleteEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.delete")); err != nil {
			return nil, err
		}

		uuid := chi.URLParam(r, "uuid")
		return struct {
			UUID string `json:"uuid"`
		}{
			uuid,
		}, m.Delete(ctx, newIdentifier(uuid))
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeUpdateEndpoint for the user service
// API-Path: /users/{uuid}
func MakeUpdateEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.update")); err != nil {
			return nil, err
		}

		uuid := chi.URLParam(r, "uuid")

		inc := NewIncomplete(0, "", false)

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return struct {
				UUID string `json:"uuid"`
			}{
				uuid,
			}, m.Update(ctx, newComplete(
				newIdentifier(uuid),
				inc,
			))
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeGetRolesEndpoint for the user service
// API-Endpoint: GET /users/{uuid}/roles
func MakeGetRolesEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.roles.get")); err != nil {
			return nil, err
		}

		uuid := chi.URLParam(r, "uuid")

		return m.GetRoles(ctx, newIdentifier(uuid))
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}

// MakeSetRolesEndpoint for the user service
// API-Endpoint: PATCH /users/{uuid}/roles
func MakeSetRolesEndpoint(l *zap.Logger, m *Manager, e encode.Encoder, rb rbac.Control, pubKey rsa.PublicKey) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if rb.IsSubjectAllowed(ctx, rbac.SubjectID(tok.Data().ID), rbac.Rule("users.roles.set")); err != nil {
			return nil, err
		}

		uuid := chi.URLParam(r, "uuid")

		roles := make(rbac.SubjectRoles, 0)
		if err := json.NewDecoder(r.Body).Decode(&roles); err != nil {
			return nil, err
		}

		return struct{}{}, m.SetRoles(ctx, newIdentifier(uuid), roles)
	}).WithBefore(token.NewMiddleware(pubKey)).HandlerFunc(l)
}
