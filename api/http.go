package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Routes - build routes
func Routes(logger *log.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	api := NewLogPerfAPI(logger)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/logperf", api.LogperfRoutes())
	})

	return router
}

// RunServer - run http server
func RunServer(logger *log.Logger) {
	router := Routes(logger)

	logger.Fatal(http.ListenAndServe(":8080", router))
}
