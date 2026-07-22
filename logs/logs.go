// Log in single-line JSON, for CloudWatch and Railway structured-log parsing alike.
// Print is the exception — dev-only, pretty-printed for terminal readability.
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
	Level    Level  `json:"level"`
	Request  any    `json:"request,omitempty"`
	Function string `json:"function,omitempty"`
	AppEnv   string `json:"app_env,omitempty"`
	Service  string `json:"service,omitempty"`
	Message  string `json:"message"`
	Data     any    `json:"data,omitempty"`
}

var Request any

// Service identifies which upstream microservice is logging, since booky-ark
// bundles many microservices' handlers into one process. Set by the caller
// (e.g. booky-ark's apigwadapter/eventadapter) before invoking a handler.
var Service string

// Production-only debugging. Suppressed outside production (APP_ENV != production).
// Use for targeted diagnostics when investigating a live issue.
func Debug(note string, data ...any) {
	if os.Getenv("APP_ENV") != PRODUCTION_ENV {
		return
	}

	logIt(DEBUG, note, data...)
}

// Dev-only pretty-printed output. Suppressed in production (APP_ENV=production).
func Print(note string, data ...any) {
	if os.Getenv("APP_ENV") == PRODUCTION_ENV {
		return
	}

	logIt(PRINT, note, data...)
}

// Business event recording — analytics-shaped breadcrumbs.
// Use for "user did X" events: payment initiated, promo applied, order placed.
func Trace(note string, data ...any) {
	logIt(TRACE, note, data...)
}

// Operational milestones and non-business parse rejections from external callers
// (bad JSON, missing fields, invalid UUIDs). Do not use for variable dumps — use Debug.
// Alarmed at high threshold.
func Info(note string, data ...any) {
	logIt(INFO, note, data...)
}

// Expected business failures: rule rejected (4xx), idempotent short-circuit (2xx with hiccup), retry recovered.
// Alarmed at moderate threshold.
func Warn(note string, data ...any) {
	logIt(WARN, note, data...)
}

// System misbehaved: sync 5xx, async task ended Failed, retries exhausted,
// programmer bug surfacing as 500, or side-effect failure after a successful primary operation.
// Alarmed at low threshold (pages during business hours).
func Error(note string, data ...any) {
	logIt(ERROR, note, data...)
}

// Service cannot run. Pair with panic() or os.Exit(1) to halt — Fatal() alone only logs.
// Alarmed at threshold 1 in every environment (pages immediately).
func Fatal(note string, data ...any) {
	logIt(FATAL, note, data...)
}

func logIt(level Level, message string, data ...any) {
	le := logEntry{
		Level:   level,
		Message: message,
	}

	if len(data) > 0 {
		le.Data = data[0]
	}

	if level != PRINT {
		le.Request = Request
		le.Function = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
		le.AppEnv = os.Getenv("APP_ENV")
		le.Service = Service
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

// Print is dev-only console output — pretty-printed for terminal readability.
// It never runs in production, so it never reaches Railway/CloudWatch. Every
// other level ships single-line JSON so line-oriented log collectors can
// parse one entry per line.
func jsonMarshal(le logEntry) ([]byte, error) {
	if le.Level == PRINT {
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
