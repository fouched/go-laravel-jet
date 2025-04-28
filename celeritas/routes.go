package celeritas

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (c *Celeritas) routes() http.Handler {
	mux := chi.NewRouter()
	addMiddleware(mux, c)

	return mux
}

func addMiddleware(mux *chi.Mux, c *Celeritas) {
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	//if c.Debug {
	//	mux.Use(middleware.Logger)
	//}

	mux.Use(c.SessionLoad)
	mux.Use(c.NoSurf)
}
