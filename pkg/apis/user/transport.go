package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/encode"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// MakeGetEndpoint for the user service
// API-Path: /users/{uuid}
func MakeGetEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		// TODO: insert scope
		return m.Get(ctx, newIdentifier(chi.URLParam(r, "uuid")))
	}).HandlerFunc(l)
}

// MakeGetByGameSerialHashEndpoint for the user service
// API-Path: /users/hash/{hash}
func MakeGetByGameSerialHashEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		// TODO: insert scope
		return m.GetByGameSerialHash(ctx, chi.URLParam(r, "hash"))
	}).HandlerFunc(l)
}

// MakeCreateEndpoint for the user service
// API-Path: /users
func MakeCreateEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		// TODO: insert scope
		inc := NewIncomplete(0, "", false)

		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return m.Create(ctx, inc)
	}).HandlerFunc(l)
}

// MakeDeleteEndpoint for the user service
// API-Path: /users/{uuid}
func MakeDeleteEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		// TODO: insert scope
		uuid := chi.URLParam(r, "uuid")
		return struct {
			UUID string `json:"uuid"`
		}{
			uuid,
		}, m.Delete(ctx, newIdentifier(uuid))
	}).HandlerFunc(l)
}

// MakeUpdateEndpoint for the user service
// API-Path: /users/{uuid}
func MakeUpdateEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		// TODO: insert scope
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
	}).HandlerFunc(l)
}
