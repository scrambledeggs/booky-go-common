# Idempotency Package

A Go package for implementing idempotent API handlers in AWS Lambda functions.

## Overview

The idempotency package provides middleware for AWS Lambda functions that handles idempotent API requests automatically. It ensures that identical requests result in identical responses, preventing duplicate operations when clients retry requests.

## Features

- Automatic detection of duplicate requests based on request body and configurable headers
- Configurable expiration time for idempotency records
- DynamoDB-based persistence layer for tracking request status
- Support for both AWS-hosted and local DynamoDB instances
- Handles in-progress requests to prevent race conditions
- Automatic cleanup of expired records using DynamoDB TTL

## Installation

```go
import "github.com/scrambledeggs/booky-go-common/idempotency"
```

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/scrambledeggs/booky-go-common/idempotency"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Your handler logic here
}

func main() {
    // Wrap your handler with the idempotency middleware
    lambda.Start(idempotency.NewIdempotentHandler(handler))
}
```

### Advanced Configuration

```go
import (
    "context"
    "time"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/scrambledeggs/booky-go-common/idempotency"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Your handler logic here
}

func main() {
    // Configure custom options
    duration := 24 * time.Hour
    tableName := "custom-idempotency-table"
    dbUrl := "http://localhost:8000" // For local development

    options := idempotency.IdempotentHandlerOptions{
        ExpiryDuration:                   &duration,
        TableName:                        &tableName,
        DynamoDBUrl:                      &dbUrl,
        HeadersToIncludeInKey:            []string{"Authorization", "X-Request-ID"},
    }

    // Wrap your handler with the configured idempotency middleware
    lambda.Start(idempotency.NewIdempotentHandlerWithOptions(handler, options))
}
```

## Configuration Options

The idempotency middleware can be configured with the following options:

| Option | Description | Default |
|--------|-------------|---------|
| `TableName` | DynamoDB table name | Value of `IDEMPOTENCY_DB_TABLE` environment variable |
| `DynamoDBUrl` | URL for DynamoDB instance | AWS DynamoDB service |
| `ExpiryDuration` | How long to store idempotency records | 1 hour |
| `HeadersToIncludeInKey` | List of HTTP header keys to include in the idempotency key hash (along with the request body) | [] (body only) |

## DynamoDB Table Structure

The idempotency package requires a DynamoDB table with the following structure:

- Primary Key: `idempotency_key` (String)
- Sort Key: `http_method#path` (String)
- Additional Attributes:
  - `response` (String): Serialized API response
  - `status` (String): Request status (in_progress, completed)
  - `expiration` (String): Expiration timestamp
  - `error` (String): Error message if applicable
  - `request_headers` (String): Original request headers
  - `ttl` (Number): Time-to-live timestamp for automatic cleanup
- Note:
  - `ttl vs expiration`: ttl is the time-to-live timestamp for automatic cleanup, expiration is the idempotency time duration

### CloudFormation Template

You can use the following CloudFormation template in your `template.yaml` file to create the required DynamoDB table with auto-scaling capabilities:

```yaml
Resources:
  IdempotencyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      KeySchema:
        - AttributeName: idempotency_key
          KeyType: HASH
        - AttributeName: http_method#path
          KeyType: RANGE
      AttributeDefinitions:
        - AttributeName: idempotency_key
          AttributeType: S
        - AttributeName: http_method#path
          AttributeType: S
      BillingMode: PAY_PER_REQUEST
      TimeToLiveSpecification:
        AttributeName: ttl
        Enabled: true

  YourLambdaFunctionRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service:
                - lambda.amazonaws.com
            Action:
              - "sts:AssumeRole"
      ManagedPolicyArns:
        - "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
      Policies:
        - PolicyName: IdempotencyDynamoDBAccess
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action:
                  - dynamodb:GetItem
                  - dynamodb:PutItem
                  - dynamodb:UpdateItem
                Resource: !GetAtt IdempotencyTable.Arn
        # Add other policies here
        # - PolicyName: SNSPublishAccess
        #   PolicyDocument:
        #     Version: '2012-10-17'
        #     Statement:
        #       - Effect: Allow
        #         Action:
        #           - sns:Publish
        #         Resource: !Ref SomeTopic

  YourLambdaFunction:
    Type: AWS::Serverless::Function
    Properties:
      # ... other properties
      Role: !GetAtt YourLambdaFunctionRole.Arn
```

## How It Works

1. When a request is received, the middleware generates an idempotency key from the request body and optionally specified headers.
2. It checks if a record with this key already exists in DynamoDB.
3. If a record exists and is still valid (not expired):
   - If the status is "completed", it returns the cached response.
   - If the status is "in_progress", it returns a 425 Too Early error.
4. If no record exists or the existing record is expired, it:
   - Creates a new record with status "in_progress".
   - Calls the handler function.
   - Updates the record with the response and status "completed".
   - Sets a TTL value for automatic cleanup (7 days after expiration).
5. DynamoDB automatically removes expired records based on the TTL attribute.

### Header Inclusion

By default, the idempotency key is generated from the request body.

## Error Handling

If the original handler returns an error, the middleware will:
1. Cache the error along with the response.
2. Return the same error for identical requests within the expiry period.

This ensures that error responses are also idempotent.

## Best Practices for Header Inclusion

To improve idempotency reliability in high-concurrency environments:

- Include headers like `Authorization` or `X-Request-ID` when retries may contain the same body but represent different logical requests.
- Avoid using headers with volatile values (e.g., `User-Agent`, timestamps) unless necessary.
- If unsure, start without headers and monitor for edge cases.

This feature gives you the flexibility to strike the right balance for your use case.

## TODO

- Add comprehensive benchmarks
- Add integration tests with DynamoDB or LocalStack
