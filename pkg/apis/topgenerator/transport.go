package top

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/encode"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// MakeHTTPGetEndpoint for tops
func MakeHTTPGetEndpoint(l *zap.Logger, e encode.Encoder, m *Manager) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		p, err := parseParams(r)
		if err != nil {
			return nil, err
		}

		return m.Get(ctx, newIdentifier(
			p.sex,
			p.undershirtID,
			p.undershirtTextureID,
			p.overshirtID,
			p.overshirtTextureID,
		))
	}).HandlerFunc(l)
}

type httpComplete struct {
	Identifier
	Incomplete
}

// MakeHTTPUpsertEndpoint for top objects
func MakeHTTPUpsertEndpoint(l *zap.Logger, e encode.Encoder, m *Manager) http.HandlerFunc {
	return endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		p, err := parseParams(r)
		if err != nil {
			return nil, err
		}

		inc := NewIncomplete("", "", 0, 0, 0, 0, 0, 25, 25, 25, 25, 1)
		if err := json.NewDecoder(r.Body).Decode(&inc); err != nil {
			return nil, err
		}

		return struct{}{}, m.Upsert(ctx, &httpComplete{
			newIdentifier(
				p.sex,
				p.undershirtID,
				p.undershirtTextureID,
				p.overshirtID,
				p.overshirtTextureID,
			),
			inc,
		})
	}).HandlerFunc(l)
}

type httpRequestIDs struct {
	sex                 uint8
	undershirtID        uint
	undershirtTextureID uint
	overshirtID         uint
	overshirtTextureID  uint
}

func parseParams(r *http.Request) (*httpRequestIDs, error) {
	sex, err := strconv.ParseUint(chi.URLParam(r, "sex"), 10, 64)
	if err != nil {
		return nil, err
	}

	undershirtID, err := strconv.ParseUint(chi.URLParam(r, "undershirtId"), 10, 64)
	if err != nil {
		return nil, err
	}

	undershirtTextureID, err := strconv.ParseUint(chi.URLParam(r, "undershirtTextureId"), 10, 64)
	if err != nil {
		return nil, err
	}

	overshirtID, err := strconv.ParseUint(chi.URLParam(r, "overshirtId"), 10, 64)
	if err != nil {
		return nil, err
	}

	overshirtTextureID, err := strconv.ParseUint(chi.URLParam(r, "overshirtTextureId"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &httpRequestIDs{
		uint8(sex),
		uint(undershirtID),
		uint(undershirtTextureID),
		uint(overshirtID),
		uint(overshirtTextureID),
	}, nil
}
