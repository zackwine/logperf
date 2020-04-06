package loggen

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// LogGenerator : Used to generate JSON like logs messages
type LogGenerator struct {
	log                 *log.Logger
	baseMap             map[string]interface{}
	timestampFieldName  string
	timestampOffsetDays int
	counterFieldname    string
	seqNum              int
	randstrgen          *RandStringGen
}

// NewLogGenerator : Initialize a LogGenerator
func NewLogGenerator(fields map[string]interface{}, timestampField string, counterField string, daysoffset int, logger *log.Logger) *LogGenerator {
	l := &LogGenerator{}

	if timestampField == "" {
		l.timestampFieldName = "@timestamp"
	} else {
		l.timestampFieldName = timestampField
	}
	if counterField == "" {
		l.timestampFieldName = "seqNum"
	} else {
		l.counterFieldname = counterField
	}
	l.timestampOffsetDays = daysoffset
	l.log = logger

	l.baseMap = make(map[string]interface{})

	// Create a regex to help generate random strings for padding
	randregex := "^<RAND:([0-9]+)>$"
	randMatcher, err := regexp.Compile(randregex)
	if err != nil {
		l.log.Println("Failed to compile regex " + randregex)
	}

	// Set default timestamp field
	l.randstrgen = NewRandStringGen()

	for key, val := range fields {

		//  If this random JSON value is a string treat it as such
		if strval, ok := val.(string); ok {
			randMatches := randMatcher.FindStringSubmatch(strval)
			if randMatches != nil {
				randStrLen, err := strconv.Atoi(randMatches[1])
				if err == nil {
					l.log.Println("Failed to convert RAND val to string " + strval)
				}
				l.baseMap[key] = l.randstrgen.RandString(randStrLen)
			} else if strval == "<UUID>" {
				uuidStr := uuid.NewV4().String()
				// Remove dashes '-' from uuid generated above
				uuidStr = strings.Replace(uuidStr, "-", "", -1)
				l.baseMap[key] = uuidStr
			} else {
				l.baseMap[key] = val
			}

		} else {
			l.baseMap[key] = val
		}
	}

	return l
}

// SetField :
func (l *LogGenerator) SetField(field string, value string) {
	l.baseMap[field] = value
}

// ResetSeqNum :
func (l *LogGenerator) ResetSeqNum() {
	l.seqNum = 0
}

// GetMessage : Generate a message
func (l *LogGenerator) GetMessage(timestamp time.Time) (string, error) {

	l.baseMap[l.counterFieldname] = strconv.Itoa(l.seqNum)
	l.seqNum++
	offsetstr := timestamp.AddDate(0, 0, -l.timestampOffsetDays).Format(time.RFC3339)
	l.baseMap[l.timestampFieldName] = offsetstr

	messagebytes, err := json.Marshal(l.baseMap)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(messagebytes[:]), err
}
