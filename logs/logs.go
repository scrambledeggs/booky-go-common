// Log in JSON format for cloudwatch
// All level entries are the same except for Debug
// Level will be used for different parameters
package logs

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws"
)

type Level string

const (
	DEBUG Level = "Debug"
	PRINT Level = "Print"
	INFO  Level = "Info"
	TRACE Level = "Trace"
	WARN  Level = "Warn"
	ERROR Level = "Error"
	FATAL Level = "Fatal"
)

type PublishFlag int

const (
	TO_SLACK PublishFlag = iota + 1
)

const PRODUCTION_ENV = "production"

type logEntry struct {
	Level    Level    `json:"level"`
	Env      []string `json:"env,omitempty"`
	Request  any      `json:"request,omitempty"`
	Function string   `json:"function,omitempty"`
	AppEnv   string   `json:"app_env,omitempty"`
	Note     string   `json:"note"`
	Data     any      `json:"data,omitempty"`
}

var Request any

// For development debugging. Will not log on Production.
// Requires `APP_ENV` equal `production`
func Debug(note string, data ...any) {
	if os.Getenv("APP_ENV") == PRODUCTION_ENV {
		return
	}

	logIt(DEBUG, note, data...)
}

// For development logging only. Will not log on Production.
// Requires `APP_ENV` equal `production`
func Print(note string, data ...any) {
	logIt(PRINT, note, data...)
}

// Log used in event tracing like User clicked pay (CT like)
func Trace(note string, data ...any) {
	logIt(TRACE, note, data...)
}

// Basic level logging
func Info(note string, data ...any) {
	logIt(INFO, note, data...)
}

// Level for expected errors like ID not found
func Warn(note string, data ...any) {
	logIt(WARN, note, data...)
}

// Level for "expected" crashes like can't connect to database
func Error(note string, data ...any) {
	logIt(ERROR, note, data...)
}

// Level for system crashes
func Fatal(note string, data ...any) {
	logIt(FATAL, note, data...)
}

func logIt(level Level, note string, data ...any) {
	le := logEntry{
		Level: level,
		Note:  note,
	}

	if len(data) > 0 {
		le.Data = data[0]
	}

	if level != PRINT {
		le.Env = os.Environ()
		le.Request = Request
		le.Function = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
		le.AppEnv = os.Getenv("APP_ENV")
	}

	l, err := jsonMarshal(le)
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		if len(data) <= 1 {
			return
		}

		for i := 1; i < len(data); i++ {
			if data[i].(PublishFlag) == TO_SLACK {
				publishToSNS(string(l))

				continue
			}
		}
	}()

	go func() {
		defer wg.Done()

		log.Print(string(l))
	}()

	wg.Wait()
}

func jsonMarshal(le logEntry) ([]byte, error) {
	if le.Level == DEBUG || le.Level == PRINT {
		return json.MarshalIndent(le, "", "  ")
	}

	return json.Marshal(le)
}

func publishToSNS(message string) {
	topicArn := os.Getenv("LOG_TO_SLACK_TOPIC_ARN")

	if topicArn == "" {
		log.Print("LOG_TO_SLACK_TOPIC_ARN env is blank")

		return
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Print("config.LoadDefaultConfig", err.Error())

		return
	}

	snsClient := sns.NewFromConfig(cfg)

	_, err = snsClient.Publish(
		ctx,
		&sns.PublishInput{
			TopicArn: aws.String(topicArn),
			Message:  aws.String(message),
		},
	)

	if err != nil {
		log.Print("snsClient.Publish", err.Error())
	}
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}
