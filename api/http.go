package api

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// HTTPServer used to create a server
type HTTPServer struct {
	Logger *log.Logger
	Addr   string
	server *http.Server
}

// routes - build routes
func (h *HTTPServer) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	api := NewLogPerfAPI(h.Logger)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/logperf", api.LogperfRoutes())
	})

	return router
}

// Start - run http server
func (h *HTTPServer) Start() {
	router := h.routes()

	h.server = &http.Server{Addr: h.Addr, Handler: router}

	go h.server.ListenAndServe()
}

// Stop - stop the http server
func (h *HTTPServer) Stop() {
	err := h.server.Shutdown(context.Background())
	if err != nil {
		h.Logger.Println("Failed to stop server")
		h.Logger.Print(err)
	}
}
