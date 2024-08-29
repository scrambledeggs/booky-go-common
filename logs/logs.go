// Log in JSON format for cloudwatch
// All level entries are the same except for Debug
// Level will be used for different parameters
package logs

import (
	"encoding/json"
	"log"
	"os"
)

const (
	DEBUG = "Debug"
	INFO  = "Info"
	TRACE = "Trace"
	WARN  = "Warn"
	ERROR = "Error"
	FATAL = "Fatal"

	PRODUCTION_ENV = "production"
)

type logEntry struct {
	Level    string   `json:"level"`
	Env      []string `json:"env"`
	Request  any      `json:"request"`
	Function string   `json:"function"`
	AppEnv   string   `json:"app_env"`
	Note     string   `json:"note"`
	Data     any      `json:"data"`
}

var Request any

// For development debugging. Will not log on Production.
// Requires `APP_ENV` equal `production`
func Debug(note string, data any) {
	if os.Getenv("APP_ENV") == PRODUCTION_ENV {
		return
	}

	logIt(DEBUG, note, data)
}

// Log used in event tracing like User clicked pay (CT like)
func Trace(note string, data any) {
	logIt(TRACE, note, data)
}

// Basic level logging
func Info(note string, data any) {
	logIt(INFO, note, data)
}

// Level for expected errors like ID not found
func Warn(note string, data any) {
	logIt(WARN, note, data)
}

// Level for "expected" crashes like can't connect to database
func Error(note string, data any) {
	logIt(ERROR, note, data)
}

// Level for system crashes
func Fatal(note string, data any) {
	logIt(FATAL, note, data)
}

func logIt(level string, note string, data any) {
	le := logEntry{
		Level:    level,
		Env:      os.Environ(),
		Request:  Request,
		Function: os.Getenv("AWS_LAMBDA_FUNCTION_NAME"),
		AppEnv:   os.Getenv("APP_ENV"),
		Note:     note,
		Data:     data,
	}

	l, err := json.Marshal(le)
	if err != nil {
		panic(err.Error())
	}

	log.Print(string(l))
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}
