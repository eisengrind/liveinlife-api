package token

import (
	"context"

	"github.com/pkg/errors"
)

// ContextKey for context payload
type ContextKey string

const (
	//TokenContextKey for API-token in a context
	TokenContextKey ContextKey = "api_bearer_token"
)

// FromContext returns a Token from a Context
func FromContext(ctx context.Context) (Token, error) {
	tok, ok := ctx.Value(TokenContextKey).(Token)
	if !ok {
		return nil, errors.New("no token in context")
	}

	return tok, nil
}

// ToContext moves a token into a Context
func ToContext(ctx context.Context, tok Token) context.Context {
	return context.WithValue(ctx, TokenContextKey, tok)
}
