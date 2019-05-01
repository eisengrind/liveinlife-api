package preselect

import (
	"net/http"

	"github.com/51st-state/api/pkg/encode"
	"go.uber.org/zap"
)

// MakeGetLeftPreSelectionsEndpoint for a HTTP server
func MakeGetLeftPreSelectionsEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return nil /*endpoint.New(e, func(ctx context.Context, _ *http.Request) (interface{}, error) {
		left, err := m.GetLeft(ctx)
		if err != nil {
			return nil, err
		}

		return struct {
			Count uint64 `json:"count"`
		}{
			left,
		}, nil
	}).HandlerFunc(l)*/
}

// MakeGetNextPreSelectionsEndpoint for a HTTP server
func MakeGetNextPreSelectionsEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return nil /*endpoint.New(e, func(ctx context.Context, _ *http.Request) (interface{}, error) {
		return m.GetNext(ctx)
	}).HandlerFunc(l)*/
}

// MakeCreatePreSelectionsEndpoint for a HTTP server
func MakeCreatePreSelectionsEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return nil /*endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		cInternal := make([]*complete, 0)

		if err := json.NewDecoder(r.Body).Decode(&cInternal); err != nil {
			return nil, err
		}

		c := make([]Complete, 0)

		for _, v := range cInternal {
			c = append(c, v)
		}

		return struct{}{}, m.Create(ctx, c...)
	}).HandlerFunc(l)*/
}

// MakeSetPreSelectionsEndpoint for a HTTP server
func MakeSetPreSelectionsEndpoint(l *zap.Logger, m *Manager, e encode.Encoder) http.HandlerFunc {
	return nil /*endpoint.New(e, func(ctx context.Context, r *http.Request) (interface{}, error) {
		cInternal := make([]*complete, 0)

		if err := json.NewDecoder(r.Body).Decode(&cInternal); err != nil {
			return nil, err
		}

		c := make([]Complete, 0)

		for _, v := range cInternal {
			c = append(c, v)
		}

		return struct{}{}, m.Set(ctx, c...)
	}).HandlerFunc(l)*/
}
