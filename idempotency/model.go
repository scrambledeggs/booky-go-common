package idempotency

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type idempotencyRecord struct {
	IdempotencyKey string `dynamodbav:"idempotency_key"`
	HttpMethodPath string `dynamodbav:"http_method#path"`
	Error          string `dynamodbav:"error"`
	Expiration     string `dynamodbav:"expiration"`
	Ttl            int64  `dynamodbav:"ttl"`
	RequestHeaders string `dynamodbav:"request_headers"`
	Response       string `dynamodbav:"response"`
	Status         string `dynamodbav:"status"`
}

var (
	idempotencyStatusInProgress = "in_progress"
	idempotencyStatusCompleted  = "completed"
	// idempotencyStatusExpired    = "expired"
)

// GetKey returns the primary key in a format that can be sent to DynamoDB.
func (ir idempotencyRecord) GetKey() map[string]types.AttributeValue {
	idempotencyKey, err := attributevalue.Marshal(ir.IdempotencyKey)
	if err != nil {
		panic(err.Error())
	}

	httpMethodPath, err := attributevalue.Marshal(ir.HttpMethodPath)
	if err != nil {
		panic(err.Error())
	}

	return map[string]types.AttributeValue{"idempotency_key": idempotencyKey, "http_method#path": httpMethodPath}
}

// Item converts the IdempotencyRecord struct to a type that can be inserted to Dynamodb via PutItem()
func (ir idempotencyRecord) Item() map[string]types.AttributeValue {
	item, err := attributevalue.MarshalMap(ir)

	if err != nil {
		panic(err.Error())
	}

	return item
}
