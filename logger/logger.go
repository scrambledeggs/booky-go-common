package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

var loggerInstance Logger

func init() {
	loggerInstance = Logger{
		infoLogger:    log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		warningLogger: log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime),
		errorLogger:   log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime),
	}
}

func Info(args ...any) {
	loggerInstance.infoLogger.Println(args...)
}

func Infof(formattedStr string, args ...any) {
	loggerInstance.infoLogger.Printf(formattedStr, args...)
}

func Warn(args ...any) {
	loggerInstance.warningLogger.Println(args...)
}

func Warnf(formattedStr string, args ...any) {
	loggerInstance.warningLogger.Printf(formattedStr, args...)
}

func Error(args ...any) {
	loggerInstance.errorLogger.Println(args...)
}

func Errorf(formattedStr string, args ...any) {
	loggerInstance.errorLogger.Printf(formattedStr, args...)
}
