// Log in single-line JSON, for CloudWatch and Railway structured-log parsing alike.
// Print is the exception — dev-only, pretty-printed for terminal readability.
// Level will be used for different parameters
package logs

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strings"
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
		caller := callerFuncName()

		le.Request = Request
		le.Function = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
		if le.Function == "" {
			// Real AWS Lambda always sets this env var; BCP/booky-ark and
			// local dev never do, so fall back to deriving it from the
			// caller's own stack frame.
			le.Function = handlerNameFromFuncName(caller)
		}
		le.AppEnv = os.Getenv("APP_ENV")
		le.Service = serviceFromFuncName(caller)
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

// packageImportPath is this package's own import path — used to walk past
// its own frames in callerFuncName.
const packageImportPath = "github.com/scrambledeggs/booky-go-common/logs"

// callerFuncName identifies who logged this entry, by walking the call stack
// to the first frame outside this package and returning its fully-qualified
// Go function name — i.e. whichever function called Debug/Info/Warn/etc. Both
// serviceFromFuncName and handlerNameFromFuncName parse this same string.
//
// It's derived fresh on every call rather than read from a shared variable
// set by some other caller, so it's correct under concurrent use: a process
// like booky-ark bundles many microservices' handlers into one binary, and a
// package-level "current caller" variable would race across concurrently
// running handlers/messages, potentially misattributing one service's log
// line to another.
func callerFuncName() string {
	var pcs [32]uintptr
	n := runtime.Callers(1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, packageImportPath+".") {
			return frame.Function
		}
		if !more {
			return ""
		}
	}
}

// serviceFromFuncName derives a service name from a fully-qualified function
// name (as returned by runtime.Frame.Function), e.g.
// "github.com/scrambledeggs/booky-athena/functions/X/handler.Handler" -> "booky-athena".
//
// Most upstream services declare "module github.com/scrambledeggs/X" in
// their own go.mod, so this returns X. A few declare a bare module name with
// no scrambledeggs/ prefix at all (e.g. "module booky-freyja") — for those,
// this falls back to the first path segment, which is already the whole
// module name.
func serviceFromFuncName(name string) string {
	parts := strings.Split(name, "/")
	for i, p := range parts {
		if p == "scrambledeggs" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	return parts[0]
}

// handlerNameFromFuncName derives the upstream Lambda function's logical
// name from a fully-qualified function name, e.g.
// "github.com/scrambledeggs/booky-locations/functions/ReverseGeocodeV1/handler.Handler" -> "ReverseGeocodeV1"
// "booky-orders/functions/GetOrderV1/handler.(*service).Handle" -> "GetOrderV1"
//
// Every upstream handler's exported entry point lives in a package literally
// named "handler" (the upstream contract this monorepo documents); the
// directory containing that package is the function's logical name — this
// mirrors what AWS_LAMBDA_FUNCTION_NAME holds in real AWS. Returns "" when
// the caller doesn't match this shape (e.g. booky-ark's own inline routes).
func handlerNameFromFuncName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) < 2 {
		return ""
	}

	if !strings.HasPrefix(parts[len(parts)-1], "handler.") {
		return ""
	}

	return parts[len(parts)-2]
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
