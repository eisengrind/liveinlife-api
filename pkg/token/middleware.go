package token

import (
	"context"
	"crypto/rsa"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/pkg/errors"
)

var errInvalidToken = errors.New("invalid token given")

// NewMiddleware of token for a http request
// Moves a token from the authorization header to
// the context of a request
func NewMiddleware(pK rsa.PublicKey) endpoint.MiddlewareFunc {
	return func(ctx context.Context, r *http.Request) (context.Context, error) {
		tokStr, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
		if err != nil {
			return ctx, err
		}

		tok, err := NewFromString(&pK, tokStr)
		if err != nil {
			return nil, err
		}

		if tok.Data().User.Type != "user" &&
			tok.Data().User.Type != "service_account" {
			return nil, errInvalidToken
		}

		return ToContext(ctx, tok), nil
	}
}
