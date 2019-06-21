package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/zackwine/logperf/api"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
)

func TestAPICreate(t *testing.T) {
	httpServer := &api.HTTPServer{Logger: logger, Addr: ":8080"}
	// Server runs async
	httpServer.Start()
	defer httpServer.Stop()

	testConf := map[string]interface{}{
		"Name":         "test1",
		"output":       "stdout",
		"count":        2,
		"period":       40,
		"routines":     1,
		"timefield":    "logtime",
		"counterfield": "seqNum",
		"fields": map[string]interface{}{
			"sessionId": "<UUID>",
			"message":   "<RAND:310>",
			"component": "logperf",
			"loglevel":  "warn",
		},
	}

	testJSON, err := json.Marshal(testConf)
	if err != nil {
		logger.Fatal(err)
	}
	testBuf := bytes.NewBuffer([]byte(testJSON))

	resp, err := http.Post("http://127.0.0.1:8080/v1/api/logperf", "json", testBuf)
	if err != nil {
		logger.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		logger.Println("Create API returned wrong status code.")
		logger.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Print(string(body))
}
