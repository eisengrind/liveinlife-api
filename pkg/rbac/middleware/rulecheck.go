package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
	"github.com/51st-state/api/pkg/rbac"
	"github.com/51st-state/api/pkg/token"
)

// NewRulecheck middleware to check whether a token has needed rules.
//
// This middleware needs the token middleware to be called before.
func NewRulecheck(ctrl rbac.Control, rule rbac.Rule) endpoint.MiddlewareFunc {
	return func(ctx context.Context, r *http.Request) (context.Context, error) {
		tok, err := token.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		allowed, err := ctrl.IsAccountAllowed(
			ctx,
			rbac.AccountID(tok.Data().User.String()),
			rule,
		)
		if err != nil {
			return nil, err
		}

		if !allowed {
			return nil, errors.New("insufficient permissions")
		}

		return nil, nil
	}
}
