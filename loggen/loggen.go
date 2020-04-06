package loggen

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	randomStrings       map[string][]string
	randomInts          map[string][]int64
	uuidStrings         map[string][]string
	enumStrings         map[string][]string
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
		l.counterFieldname = "seqNum"
	} else {
		l.counterFieldname = counterField
	}
	l.timestampOffsetDays = daysoffset
	l.log = logger

	l.baseMap = make(map[string]interface{})
	l.randomStrings = make(map[string][]string)
	l.randomInts = make(map[string][]int64)
	l.uuidStrings = make(map[string][]string)
	l.enumStrings = make(map[string][]string)

	// Set some defaults for require fields
	l.baseMap["component"] = "logperf"
	l.baseMap["action"] = "blast"
	l.baseMap["loglevel"] = "warn"

	// Create a regex to help generate random strings for padding
	randregex := "^<RAND:([0-9]+)(:([0-9]+))?>$"
	randMatcher, err := regexp.Compile(randregex)
	if err != nil {
		l.log.Println("Failed to compile regex " + randregex)
	}

	// Create a regex to help generate random ints for padding
	randIntRegex := "^<RANDINT:([0-9]+)(:([0-9]+))?>$"
	randIntMatcher, err := regexp.Compile(randIntRegex)
	if err != nil {
		l.log.Println("Failed to compile regex " + randIntRegex)
	}

	// Create regex to match UUID
	UUIDRegex := "^<UUID(:([0-9]+))?>$"
	UUIDMatcher, err := regexp.Compile(UUIDRegex)
	if err != nil {
		l.log.Println("Failed to compile regex " + UUIDRegex)
	}

	// Create regex to match ENUM
	EnumRegex := "^<ENUM:(.*)>$"
	EnumMatcher, err := regexp.Compile(EnumRegex)
	if err != nil {
		l.log.Println("Failed to compile regex " + EnumRegex)
	}

	l.randstrgen = NewRandStringGen()

	for key, val := range fields {

		//  If this is a string
		if strval, ok := val.(string); ok {
			randMatches := randMatcher.FindStringSubmatch(strval)
			randIntsMatches := randIntMatcher.FindStringSubmatch(strval)
			UUIDMatches := UUIDMatcher.FindStringSubmatch(strval)
			EnumMatches := EnumMatcher.FindStringSubmatch(strval)

			if randMatches != nil {
				randStrLen, err := strconv.Atoi(randMatches[1])
				if err != nil {
					l.log.Printf("Failed to convert RAND length from string (%s) %v\n", strval, err)
					randStrLen = 10
				}

				randomCnt, err := strconv.Atoi(randMatches[3])
				if err != nil {
					if len(randMatches[3]) != 0 {
						l.log.Printf("Failed to convert RAND count from string (%s) %v\n", strval, err)
					}
					randomCnt = 1
				}
				l.generateRandomStrings(key, randStrLen, randomCnt)

			} else if randIntsMatches != nil {
				randIntLen, err := strconv.Atoi(randIntsMatches[1])
				if err != nil {
					l.log.Printf("Failed to convert RAND length from string (%s) %v\n", strval, err)
					randIntLen = 8
				}

				randomCnt, err := strconv.Atoi(randIntsMatches[3])
				if err != nil {
					if len(randIntsMatches[3]) != 0 {
						l.log.Printf("Failed to convert RAND count from string (%s) %v\n", strval, err)
					}
					randomCnt = 1
				}
				l.generateRandomInts(key, randIntLen, randomCnt)

			} else if UUIDMatches != nil {

				UUIDCnt, err := strconv.Atoi(UUIDMatches[2])
				if err != nil {
					if len(UUIDMatches[2]) != 0 {
						l.log.Printf("Failed to convert UUID count from string (%s) %v\n ", UUIDMatches[2], err)
					}
					UUIDCnt = 1
				}
				l.generateUUIDs(key, UUIDCnt)

			} else if EnumMatches != nil {
				l.enumStrings[key] = strings.Split(EnumMatches[1], ",")
			} else {
				l.baseMap[key] = val
			}

		} else {
			l.baseMap[key] = val
		}
	}

	return l
}

func (l *LogGenerator) generateRandomStrings(field string, length int, count int) {
	genStrings := make([]string, count, count)
	for i := 0; i < count; i++ {
		genStrings[i] = l.randstrgen.RandString(length)
	}
	l.randomStrings[field] = genStrings
}

func (l *LogGenerator) generateRandomInts(field string, length int, count int) {
	genInts := make([]int64, count, count)
	maxNum := math.Pow10(length+1) - 1
	for i := 0; i < count; i++ {
		genInts[i] = l.randstrgen.RandNum(int64(maxNum))
	}
	l.randomInts[field] = genInts
}

func (l *LogGenerator) generateUUIDs(field string, count int) {
	genUUIDs := make([]string, count, count)
	for i := 0; i < count; i++ {
		// Remove dashes '-' from uuid
		genUUIDs[i] = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	}
	l.uuidStrings[field] = genUUIDs
}

// SetField :
func (l *LogGenerator) SetField(field string, value string) {
	l.baseMap[field] = value
}

// ResetSeqNum :
func (l *LogGenerator) ResetSeqNum() {
	l.seqNum = 0
}

func (l *LogGenerator) updateMessage(timestamp time.Time) {

	l.baseMap[l.counterFieldname] = strconv.Itoa(l.seqNum)
	l.seqNum++
	offsetstr := timestamp.AddDate(0, 0, -l.timestampOffsetDays).Format(time.RFC3339)
	l.baseMap[l.timestampFieldName] = offsetstr

	for key, values := range l.randomStrings {
		r := l.randstrgen.RandNum(int64(len(values)))
		l.baseMap[key] = l.randomStrings[key][r]
	}

	for key, values := range l.randomInts {
		r := l.randstrgen.RandNum(int64(len(values)))
		l.baseMap[key] = l.randomInts[key][r]
	}

	for key, values := range l.uuidStrings {
		r := l.randstrgen.RandNum(int64(len(values)))
		l.baseMap[key] = l.uuidStrings[key][r]
	}

	for key, values := range l.enumStrings {
		r := l.randstrgen.RandNum(int64(len(values)))
		l.baseMap[key] = l.enumStrings[key][r]
	}

}

// GetJSONMessage : Generate a message in JSON format
func (l *LogGenerator) GetJSONMessage(timestamp time.Time) (string, error) {

	l.updateMessage(timestamp)

	messagebytes, err := json.Marshal(l.baseMap)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(messagebytes[:]), err
}

// GetStandardMessage : Generate a message in LoggingStandard format
func (l *LogGenerator) GetStandardMessage(timestamp time.Time) (string, error) {

	var sb strings.Builder

	l.updateMessage(timestamp)

	fmt.Fprintf(&sb, "%v, %s, %s, %s", l.baseMap[l.timestampFieldName], l.baseMap["component"], l.baseMap["action"], l.baseMap["loglevel"])

	for key, value := range l.baseMap {
		if strval, ok := value.(string); ok {
			if key == l.timestampFieldName {
				continue
			}
			if key == "component" {
				continue
			}
			if key == "action" {
				continue
			}
			if key == "loglevel" {
				continue
			}

			fmt.Fprintf(&sb, ", %s=\"%s\"", key, strval)
		} else if intval, ok := value.(int64); ok {
			fmt.Fprintf(&sb, ", %s=%d", key, intval)
		} else if intval, ok := value.(int); ok {
			fmt.Fprintf(&sb, ", %s=%d", key, intval)
		}
	}

	return sb.String(), nil
}
