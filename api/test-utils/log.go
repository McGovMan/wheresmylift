package test

import (
	"encoding/json"
	"strings"

	"github.com/nsf/jsondiff"
)

type LogSink struct {
	Logs []string
}

func (l *LogSink) Write(p []byte) (n int, err error) {
	log := strings.Trim(string(p), "\n")
	l.Logs = append(l.Logs, log)

	return len(p), nil
}

func (l *LogSink) Reset() {
	l.Logs = []string{}
}

func (l *LogSink) ContainsLog(expectedLog map[string]interface{}) bool {
	opts := jsondiff.DefaultJSONOptions()
	expectedLogString, err := json.Marshal(expectedLog)
	if err != nil {
		return false
	}

	for _, log := range l.Logs {
		diff, _ := jsondiff.Compare([]byte(log), expectedLogString, &opts)
		if diff == jsondiff.FullMatch {
			return true
		}
	}

	return false
}
