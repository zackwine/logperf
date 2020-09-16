package logperf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config : Configruation to run log performance tests
type Config struct {
	// Name of the test
	Name string `json:"name"`
	// Routines - number of go routines (threads) to run concurrently
	Routines int `json:"routines"`
	// Output - The output to use 'tcp', 'http', or 'stdout'
	Output string `json:"output"`
	// OutputFormat - The outputformat to use 'JSON', or 'LoggingStandard'
	OutputFormat string `json:"outputformat"`
	// Addr - The address to send logs only applies to certain outputs
	Addr string `json:"addr"`
	// Count - The number of logs to send
	Count int64 `json:"count"`
	// Period - The period to wait between each log in milliseconds
	Period int `json:"period"`
	// IndexBase - The only applies to the elasticsearch output index base name
	IndexBase string `json:"indexbase"`

	// Daysoffset used to send logs with older timestamps
	Daysoffset   int                    `json:"daysoffset"`
	Timefield    string                 `json:"timefield"`
	CounterField string                 `json:"counterfield"`
	Fields       map[string]interface{} `json:"fields"`
}

// Configs : A list of logperf configs that can be Marshalled from yaml
type Configs struct {
	Configs []Config `yaml:"logperfs"`
}

// NewConfigs : Unmarshal LogPerfConfig from YAML file
func NewConfigs(filename string) (*Configs, error) {
	t := &Configs{}
	source, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, t)
	if err != nil {
		fmt.Printf("error: %v", err)
		panic(err)
	}
	return t, err
}
