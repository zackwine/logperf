package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/winez/logperf"
)

// LogPerfAPI -
type LogPerfAPI struct {
	log      *log.Logger
	logperfs []*logperf.LogPerf
}

// NewLogPerfAPI - Instantiate a log perf API instance
func NewLogPerfAPI(logger *log.Logger) *LogPerfAPI {
	return &LogPerfAPI{
		log: logger,
	}
}

// LogperfRoutes - Generate routes for this API
func (l *LogPerfAPI) LogperfRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{logperfID}", l.GetLogPerf)
	router.Delete("/{logperfID}", l.DeleteLogPerf)
	router.Post("/", l.CreateLogPerf)
	router.Get("/", l.GetAllLogPerfs)
	return router
}

// GetLogPerf - list information about logperfs
func (l *LogPerfAPI) GetLogPerf(w http.ResponseWriter, r *http.Request) {
	logperfID := chi.URLParam(r, "logperfID")
	logperf := logperf.Config{
		Name:   logperfID,
		Output: "stdout",
	}
	render.JSON(w, r, logperf)
}

// DeleteLogPerf - Delete a logperf test
func (l *LogPerfAPI) DeleteLogPerf(w http.ResponseWriter, r *http.Request) {
	logperfID := chi.URLParam(r, "logperfID")
	response := make(map[string]string)
	response["message"] = "Deleted logperf test " + logperfID
	render.JSON(w, r, response)
}

// CreateLogPerf - Delete a logperf test
func (l *LogPerfAPI) CreateLogPerf(w http.ResponseWriter, r *http.Request) {
	var cfg logperf.Config
	response := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cfg)
	if err != nil {
		l.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		response["message"] = "Failed to decode JSON"
		response["error"] = err.Error()

	} else {
		l.log.Print(cfg)
		response["message"] = "Created logperf test"
		curPerf := logperf.NewLogPerf(cfg, l.log)

		// TODO validate JSON further
		// TODO ensure unique test names
		l.logperfs = append(l.logperfs, curPerf)
		go func(lp *logperf.LogPerf) {
			err := lp.Start(nil)
			if err != nil {
				l.log.Println(err)
			}
		}(curPerf)
	}

	render.JSON(w, r, response)
}

// GetAllLogPerfs - Delete a logperf test
func (l *LogPerfAPI) GetAllLogPerfs(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "All logperf tests"
	render.JSON(w, r, response)
}
