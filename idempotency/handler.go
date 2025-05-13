package idempotency

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	r "github.com/scrambledeggs/booky-go-common/apigatewayresponse"
	"github.com/scrambledeggs/booky-go-common/logs"
)

var timeFormat = time.RFC3339

var (
	defaultExpiryTimeDuration = 1 * time.Hour
	defaultDynamoDBTableName  = os.Getenv("IDEMPOTENCY_DB_TABLE")
)

type lambdaHandlerFunc func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type IdempotentHandlerOptions struct {
	DynamoDBUrl    *string
	TableName      *string
	ExpiryDuration *time.Duration
}

type idempotentHandler struct {
	lambdaHandler  lambdaHandlerFunc
	tableName      string
	dynamoDBUrl    *string
	expiryDuration time.Duration
}

// NewIdempotentHandler attaches a middleware to a lambda handler with default DynamoDB persistence store to enable idempotency in each response
func NewIdempotentHandler(handler lambdaHandlerFunc) lambdaHandlerFunc {
	return idempotentHandler{lambdaHandler: handler, tableName: defaultDynamoDBTableName, expiryDuration: defaultExpiryTimeDuration}.handler
}

// NewIdempotentHandlerWithOptions attaches a middleware to a lambda handler with configurable database options to enable idempotency in each response.

// Parameters:
//   - `TableName` - required: no, default: env for `IDEMPOTENCY_DB_TABLE`. The dynamo db table name to execute queries
//   - `DynamoDBUrl` - required: no, default: nil. Optional parameter than can be used when working with local DynamoDB instance
//   - `ExpiryDuration` - required: no, default: 1 hour. Specify expiration per request before processing another non-idempotent response
func NewIdempotentHandlerWithOptions(handler lambdaHandlerFunc, options IdempotentHandlerOptions) lambdaHandlerFunc {
	expiryDuration := &defaultExpiryTimeDuration
	tableName := defaultDynamoDBTableName

	if options.ExpiryDuration != nil {
		expiryDuration = options.ExpiryDuration
	}

	if options.TableName != nil {
		tableName = *options.TableName
	}

	return idempotentHandler{lambdaHandler: handler, tableName: tableName, dynamoDBUrl: options.DynamoDBUrl, expiryDuration: *expiryDuration}.handler
}

func (ih idempotentHandler) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idempotencyKey := mapToHash(request.Body)
	httpMethodPath := fmt.Sprintf("%s#%s", request.HTTPMethod, request.Path)
	requestHeaders, _ := json.Marshal(request.Headers)

	logs.Print("data", map[string]any{
		"key":  idempotencyKey,
		"body": request.Body,
	})

	dbClient := newIdempotencyDBClient(context.TODO(), ih.tableName)
	if ih.dynamoDBUrl != nil {
		dbClient = newIdempotencyDBClientWithUrl(context.TODO(), ih.tableName, *ih.dynamoDBUrl)
	}

	record := dbClient.Get(idempotencyRecord{IdempotencyKey: idempotencyKey, HttpMethodPath: httpMethodPath})

	if record != nil {
		expiry, _ := time.Parse(timeFormat, record.Expiration)
		currentTime := time.Now()

		if currentTime.Before(expiry) && record.Status == idempotencyStatusCompleted {
			var resp events.APIGatewayProxyResponse
			json.Unmarshal([]byte(record.Response), &resp)
			if record.Error != "" {
				return resp, errors.New(record.Error)
			}
			logs.Print("returned IDEMPOTENT response", nil)
			return resp, nil
		}

		if currentTime.Before(expiry) && record.Status == idempotencyStatusInProgress {
			logs.Print("returned IDEMPOTENCY_ALREADY_IN_PROGRESS", nil)
			return r.SingleErrorResponse(http.StatusTooEarly, r.ErrorResponseBody{
				Code:    "IDEMPOTENCY_ALREADY_IN_PROGRESS",
				Message: "An identical request is already in progress",
			})
		}
	}

	expiry := time.Now().Add(ih.expiryDuration)
	ttl := expiry.Add(7 * 24 * time.Hour).Unix()

	record = &idempotencyRecord{
		IdempotencyKey: idempotencyKey,
		HttpMethodPath: httpMethodPath,
		Expiration:     expiry.Format(timeFormat),
		Ttl:            ttl,
		RequestHeaders: string(requestHeaders),
		Status:         idempotencyStatusInProgress,
	}
	dbClient.Put(*record)

	resp, err := ih.lambdaHandler(request)
	r, _ := json.Marshal(resp)

	record.Response = string(r)
	if err != nil {
		record.Error = string(err.Error())
	}
	record.Status = idempotencyStatusCompleted
	dbClient.Put(*record)

	logs.Print("returned NON IDEMPOTENT response", nil)
	return resp, err
}

func mapToHash(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	bytes := hash.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}
