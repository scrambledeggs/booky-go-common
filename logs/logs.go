package logs

import (
	"encoding/json"
	"log"
	"os"
)

const (
	DEBUG = "Debug"
	INFO  = "Info"
	WARN  = "Warn"
	ERROR = "Error"
	FATAL = "Fatal"
)

type logEntry struct {
	Level   string   `json:"level"`
	Env     []string `json:"environment"`
	Request any      `json:"request"`
	Message string   `json:"message"`
	Data    any      `json:"data"`
}

var Request any

func Debug(message string, data any) {
	logIt(DEBUG, message, data)
}

func Info(message string, data any) {
	logIt(INFO, message, data)
}

func Warn(message string, data any) {
	logIt(WARN, message, data)
}

func Error(message string, data any) {
	logIt(ERROR, message, data)
}

func logIt(level string, message string, data any) {
	le := logEntry{
		Level:   level,
		Request: Request,
		Env:     os.Environ(),
		Message: message,
		Data:    data,
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
