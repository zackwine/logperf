package outputs

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

const (
	// IndexPeriodHourly create indexes per hour in elasticsearch
	IndexPeriodHourly = "hourly"
	// IndexPeriodDaily create indexes per daily in elasticsearch
	IndexPeriodDaily = "daily"
)

const (
	defaultIndexPeriod   = IndexPeriodHourly
	defaultDocType       = "logs"
	defaultIndexBase     = "logperf"
	esRetryCount         = 5
	esRetryBackoffFactor = 2
	// The actual AWS ES-AAS bulk request limit is 10MB, but leave a 10% buffer
	// This is the smallest limit.  Larger ES deployments have a larger limit of 100MB.
	// TODO:  Make this dynamic
	//bulkRequestLimit  int64 = 9 * 1024 * 1024
	bulkRequestLimit  int64 = 1 * 1024 * 1024 // Break it up into 1MB chunks
	indexActionString       = "{\"index\":{}}\n"
)

// ElasticSearchOutputConfig - Config for the elasticsearch output for JSON formated data
type ElasticSearchOutputConfig struct {
	Addresses    []string
	IndexBase    string
	IndexPeriod  string
	DocumentType string
}

// ElasticSearchOutput : a channel based output for perftesting logs to elasticsearch.
type ElasticSearchOutput struct {
	ElasticSearchOutputConfig
	client        *elasticsearch.Client
	log           *log.Logger
	inputchan     chan string
	stopchan      chan bool
	running       bool
	indexedCount  int64
	messageBuffer strings.Builder
	bufferLineCnt int64
}

// NewElasticSearchOutput : Initialize ElasticSearchOutput
func NewElasticSearchOutput(config ElasticSearchOutputConfig, logger *log.Logger) *ElasticSearchOutput {

	cfg := elasticsearch.Config{
		Addresses: config.Addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Fatalf("Error creating the client: %s", err)
	}

	// Make a simple call to verify ES returns a valid response
	res, err := es.Info()
	if err != nil {
		logger.Println(res)
		logger.Fatalf("Error getting response from elasticsearch (%v): %s", cfg.Addresses, err)
	}

	// Configure defaults
	if config.IndexPeriod == "" {
		config.IndexPeriod = defaultIndexPeriod
	}
	if config.IndexBase == "" {
		config.IndexBase = defaultIndexBase
	}
	if config.DocumentType == "" {
		config.DocumentType = defaultDocType
	}

	e := &ElasticSearchOutput{
		ElasticSearchOutputConfig: config,
		client:                    es,
		log:                       logger,
		stopchan:                  make(chan bool),
	}
	return e
}

// StartOutput : Implement the Output interface.
// Start reading from the channel and sending to the output.
func (e *ElasticSearchOutput) StartOutput(input chan string) error {
	e.log.Println("Starting elasticsearch output")
	e.inputchan = input

	if e.running {
		e.log.Println("This output is already running")
		return errors.New("This output is already running")
	}
	e.running = true

	go e.run()

	return nil
}

// StopOutput : Implement the Output interface
// Stop reading from the channel and sending to the output.
func (e *ElasticSearchOutput) StopOutput() error {

	e.log.Println("Stopping elasticsearch output")

	if !e.running {
		e.log.Println("This output is NOT running")
		return errors.New("This output is NOT running")
	}
	now := time.Now().UTC()
	err := e.postBulkIndex(e.messageBuffer, &now, e.bufferLineCnt)
	if err != nil {
		e.log.Fatalf("Failed to bulk post index: %v", err)
	} else {
		atomic.AddInt64(&e.indexedCount, int64(e.bufferLineCnt))
		e.bufferLineCnt = 0
		e.messageBuffer.Reset()
	}

	go func() {
		e.stopchan <- true
	}()

	return nil
}

func (e *ElasticSearchOutput) run() {

	for {
		select {
		case message := <-e.inputchan:
			//fmt.Println(message)
			lineLen := int64(len(message))
			if lineLen > bulkRequestLimit {
				e.log.Printf("ERROR: The logline size (%d) is larger than the limit (%d)\n", lineLen, bulkRequestLimit)
				continue
			}
			e.bufferLineCnt++

			if int64(e.messageBuffer.Len())+lineLen > bulkRequestLimit {
				now := time.Now().UTC()
				err := e.postBulkIndex(e.messageBuffer, &now, e.bufferLineCnt)
				if err != nil {
					e.log.Fatalf("Failed to bulk post index: %v", err)
				} else {
					atomic.AddInt64(&e.indexedCount, e.bufferLineCnt)
					e.bufferLineCnt = 0
					e.messageBuffer.Reset()
				}

			}

			_, err := e.messageBuffer.WriteString(indexActionString)
			if err != nil {
				e.log.Fatalf("Failed build string: %v", err)
			}

			_, err = e.messageBuffer.WriteString(message + "\n")
			if err != nil {
				e.log.Fatalf("Failed build string: %v", err)
			}

		case <-e.stopchan:
			e.log.Println("ElasticSearchOutput thread stopping")
			e.running = false

			close(e.stopchan)
			return
		}
	}
}

func (e *ElasticSearchOutput) getIndexName(t *time.Time) string {
	var indexSuffixFormat string
	switch e.IndexPeriod {
	case IndexPeriodDaily:
		indexSuffixFormat = "-2006-01-02"
	case IndexPeriodHourly:
		indexSuffixFormat = "-2006-01-02-15"
	}

	return e.IndexBase + t.Format(indexSuffixFormat)
}

func (e *ElasticSearchOutput) postBulkIndex(strBuilder strings.Builder, batchTime *time.Time, docCount int64) error {

	indexName := e.getIndexName(batchTime)
	e.log.Printf("Elasticsearch bulk indexing (%d) docs (%d) bytes into (%s)\n", docCount, strBuilder.Len(), indexName)
	// Set up the request object.
	req := esapi.BulkRequest{
		Index:        indexName,
		Body:         strings.NewReader(strBuilder.String()),
		DocumentType: e.DocumentType,
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		e.log.Printf("Error getting response: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		e.log.Printf("[%s] Error bulk indexing", res.Status())
		return errors.New("[" + string(res.Status()) + "] Error bulk indexing")
	}

	return nil
}

// GetCount get the number of objects indexed
func (e *ElasticSearchOutput) GetCount() int64 {
	return e.indexedCount
}
