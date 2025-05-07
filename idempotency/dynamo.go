package idempotency

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

// NewDynamoDBClient creates a new instance of DynamoDB client with default config
//
// Parameters:
//   - `dbUrl` mainly used for local debugging when passed with a local DynamoDB URL(i.e. "http://localhost:8000")
func NewDynamoDBClient(ctx context.Context, dbUrl *string) *dynamodb.Client {
	region := os.Getenv("AMZ_REGION")

	// Load default configuration with region if specified
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		panic(err.Error())
	}

	// Create and return the DynamoDB client
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if dbUrl != nil {
			o.BaseEndpoint = aws.String(*dbUrl)
		}
	})
}
