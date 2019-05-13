package auth

import (
	"context"
	"net/http"

	"github.com/51st-state/api/pkg/api/endpoint"
)

type recaptchaContext string

const recaptchaRespCtxKey recaptchaContext = "recaptcha_response"

func populateRecaptchaResponseToken() endpoint.MiddlewareFunc {
	return func(ctx context.Context, r *http.Request) (context.Context, error) {
		return recaptchaRespToCtx(ctx, r.Header.Get("X-Recaptcha-Response-Token")), nil
	}
}

func recaptchaRespToCtx(ctx context.Context, response string) context.Context {
	return context.WithValue(ctx, recaptchaRespCtxKey, response)
}

func recaptchaRespFromCtx(ctx context.Context) string {
	return ctx.Value(recaptchaRespCtxKey).(string)
}
