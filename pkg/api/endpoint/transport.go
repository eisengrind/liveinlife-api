package endpoint

import (
	"context"
	"net/http"

	"github.com/51st-state/api/pkg/encode"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/51st-state/api/pkg/problems"
	"google.golang.org/grpc/status"
)

// HandlerFunc for internal transport
type HandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

// MiddlewareFunc for internal middleware transport
type MiddlewareFunc func(ctx context.Context, r *http.Request) (context.Context, error)

// Endpoint for a http service
type Endpoint interface {
	WithBefore(mw MiddlewareFunc) Endpoint
	WithAfter(mw MiddlewareFunc) Endpoint
	HandlerFunc(*zap.Logger) http.HandlerFunc
}

type endpoint struct {
	before MiddlewareFunc
	after  MiddlewareFunc
	hndl   HandlerFunc
	enc    encode.Encoder
}

//New http endpoint for a http service
func New(e encode.Encoder, hndl HandlerFunc) Endpoint {
	return &endpoint{
		hndl: hndl,
		enc:  e,
	}
}

func (e *endpoint) WithBefore(mw MiddlewareFunc) Endpoint {
	if e.before == nil {
		e.before = mw
	} else {
		tmpBefore := e.before
		e.before = func(ctx context.Context, r *http.Request) (context.Context, error) {
			ctx, err := tmpBefore(ctx, r)
			if err != nil {
				return ctx, err
			}

			return mw(ctx, r)
		}
	}

	return e
}

func (e *endpoint) WithAfter(mw MiddlewareFunc) Endpoint {
	if e.after == nil {
		e.after = mw
	} else {
		tmpAfter := e.after
		e.after = func(ctx context.Context, r *http.Request) (context.Context, error) {
			ctx, err := tmpAfter(ctx, r)
			if err != nil {
				return ctx, errors.WithStack(err)
			}

			return mw(ctx, r)
		}
	}

	return e
}

func (e *endpoint) HandlerFunc(l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := r.Context()

		if e.before != nil {
			ctx, err = e.before(ctx, r)
			if err != nil {
				l.Error(
					"before middleware error",
					zap.Error(err),
				)
				e.encodeError(w, err, l)
				return
			}
		}

		v, err := e.hndl(ctx, r)
		if err != nil {
			l.Error(
				"handler error",
				zap.Error(err),
			)
			e.encodeError(w, err, l)
			return
		}

		if e.after != nil {
			ctx, err = e.after(ctx, r)
			if err != nil {
				l.Error(
					"after middleware error",
					zap.Error(err),
				)
				e.encodeError(w, err, l)
				return
			}
		}

		e.enc.Encode(w, v)
	}
}

func problemError(err error) *problems.Problem {
	problem, ok := err.(*problems.Problem)
	if !ok {
		st, ok := status.FromError(err)
		if ok {
			err = errors.New(st.Message())
		}

		problem = problems.New(
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			http.StatusInternalServerError,
		)
	}

	if problem.Status == 0 {
		problem.Status = http.StatusInternalServerError
	}

	return problem
}

func (e *endpoint) encodeError(w http.ResponseWriter, err error, l *zap.Logger) {
	p := problemError(err)
	w.WriteHeader(p.Status)

	if err := e.enc.Encode(w, p); err != nil {
		l.Error(
			"encoding failed",
			zap.Error(err),
		)
	}
}
