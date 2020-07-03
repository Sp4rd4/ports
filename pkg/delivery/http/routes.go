package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	l "github.com/treastech/logger"
)

func (pc *PortController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer, middleware.RequestID, l.Logger(pc.logger))
	r.Route("/ports", func(r chi.Router) {
		r.Get("/{portID}", pc.Get)
	})

	r.ServeHTTP(w, req)
}
