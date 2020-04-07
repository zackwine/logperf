package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	api := NewLogPerfAPI(h.Logger)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/logperf", api.LogperfRoutes())
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))
	fileServer(router, "/app", filesDir)

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

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
