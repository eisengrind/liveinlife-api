package api

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type loggerMiddleware struct {
	handler http.Handler
	logger  *zap.Logger
}

func (l *loggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := l.logger.With(
		zap.String("addr", r.RemoteAddr),
		zap.Strings("X-Forwarded-For", strings.Split(r.Header.Get("X-Forwarded-For"), ",")),
		zap.String("X-Real-IP", r.Header.Get("X-Real-IP")),
		zap.String("User-Agent", r.Header.Get("User-Agent")),
		zap.String("protocol", r.Proto),
		zap.String("method", r.Method),
		zap.String("uri", r.URL.RequestURI()),
		zap.String("url", r.URL.String()),
	)

	logger.Info("incoming http request")

	t := time.Now()

	defer func() {
		logger.Info(
			"http request finished",
			zap.Duration("elapsed", time.Since(t)),
		)
	}()

	l.handler.ServeHTTP(w, r)
}

func newLoggerMiddleware(h http.Handler, l *zap.Logger) http.Handler {
	return &loggerMiddleware{
		handler: h,
		logger:  l,
	}
}
