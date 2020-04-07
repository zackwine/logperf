package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"
	"github.com/zackwine/logperf"
)

// LogPerfAPI -
type LogPerfAPI struct {
	log      *log.Logger
	logperfs map[string]*logperf.LogPerf
}

// NewLogPerfAPI - Instantiate a log perf API instance
func NewLogPerfAPI(logger *log.Logger) *LogPerfAPI {
	logperfs := make(map[string]*logperf.LogPerf)
	return &LogPerfAPI{
		log:      logger,
		logperfs: logperfs,
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
	response := make(map[string]interface{})

	if lp, hasname := l.logperfs[logperfID]; hasname {
		response["id"] = logperfID
		response["logperf"] = lp
		response["target"] = lp.GetTargetCount()
		response["count"] = lp.GetCurrentCount()
		response["percentage"] = float64(lp.GetCurrentCount()) / float64(lp.GetTargetCount()) * 100
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response["error"] = "Failed to find ID (" + logperfID + ")"
	}
	render.JSON(w, r, response)
}

// DeleteLogPerf - Delete a logperf test
func (l *LogPerfAPI) DeleteLogPerf(w http.ResponseWriter, r *http.Request) {
	logperfID := chi.URLParam(r, "logperfID")
	response := make(map[string]interface{})

	if lp, hasname := l.logperfs[logperfID]; hasname {
		lp.Stop()
		delete(l.logperfs, logperfID)
		response["message"] = "Deleted logperf test " + logperfID
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response["error"] = "Failed to find ID (" + logperfID + ")"
	}

	render.JSON(w, r, response)
}

// CreateLogPerf - Create a log perf test
func (l *LogPerfAPI) CreateLogPerf(w http.ResponseWriter, r *http.Request) {
	var cfg logperf.Config

	// Extra debug for debugging decoder issues.
	/*
		requestDump, err2 := httputil.DumpRequest(r, true)
		if err2 != nil {
			l.log.Printf("Failed to dump request (%v)", err)
		}
		l.log.Printf("Received body (%s)", requestDump)
	*/

	response := make(map[string]interface{})

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cfg)
	if err != nil {
		l.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		response["message"] = "Failed to decode JSON"
		response["error"] = err.Error()
	} else {
		l.log.Print(cfg)
		curPerf := logperf.NewLogPerf(cfg, l.log)
		// Generate unique hash ID for this test
		logperfID := uuid.NewV4().String()
		// Remove dashes '-' from uuid generated above
		logperfID = strings.Replace(logperfID, "-", "", -1)

		l.logperfs[logperfID] = curPerf
		response["id"] = logperfID

		err := curPerf.Start(nil)
		if err != nil {
			l.log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			response["error"] = err
			response["message"] = "Failed to create logperf test (" + cfg.Name + ")"
		} else {
			response["message"] = "Created logperf test (" + cfg.Name + ")"
		}
	}

	render.JSON(w, r, response)
}

// GetAllLogPerfs - Delete a logperf test
func (l *LogPerfAPI) GetAllLogPerfs(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["logperfs"] = l.logperfs
	render.JSON(w, r, response)
}
