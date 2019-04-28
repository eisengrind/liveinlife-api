package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

//API object of a service
type API struct {
	chi.Router
	logger *zap.Logger
	addr   string
}

//New http api for a service
func New(addr string, l *zap.Logger) *API {
	r := chi.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return newLoggerMiddleware(h, l)
	})
	r.Use(cors.AllowAll().Handler)

	return &API{
		r,
		l,
		addr,
	}
}

//Serve a service api
func (a *API) Serve() error {
	return http.ListenAndServe(a.addr, a)
}
