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

func Debug(note string, data any) {
	if os.Getenv("APP_ENV") == PRODUCTION_ENV {
		return
	}

	logIt(DEBUG, note, data)
}

func Trace(note string, data any) {
	logIt(TRACE, note, data)
}

func Info(note string, data any) {
	logIt(INFO, note, data)
}

func Warn(note string, data any) {
	logIt(WARN, note, data)
}

func Error(note string, data any) {
	logIt(ERROR, note, data)
}

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
