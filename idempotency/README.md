# Idempotency Package

A Go package for implementing idempotent API handlers in AWS Lambda functions.

## Overview

The idempotency package provides middleware for AWS Lambda functions that handles idempotent API requests automatically. It ensures that identical requests result in identical responses, preventing duplicate operations when clients retry requests.

## Features

- Automatic detection of duplicate requests based on request body
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
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/scrambledeggs/booky-go-common/idempotency"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
    "time"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/scrambledeggs/booky-go-common/idempotency"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Your handler logic here
}

func main() {
    // Configure custom options
    duration := 24 * time.Hour
    tableName := "custom-idempotency-table"
    dbUrl := "http://localhost:8000" // For local development

    options := idempotency.IdempotentHandlerOptions{
        ExpiryDuration: &duration,
        TableName:      &tableName,
        DynamoDBUrl:    &dbUrl,
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

### CloudFormation Template

You can use the following CloudFormation template in your `template.yaml` file to create the required DynamoDB table with auto-scaling capabilities:

```yaml
Resources:
  IdempotencyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      KeySchema:
        - AttributeName: "idempotency_key"
          KeyType: "HASH"
        - AttributeName: "http_method#path"
          KeyType: "RANGE"
      AttributeDefinitions:
        - AttributeName: "idempotency_key"
          AttributeType: "S"
        - AttributeName: "http_method#path"
          AttributeType: "S"
      ProvisionedThroughput:
        ReadCapacityUnits: 5
        WriteCapacityUnits: 5
      TimeToLiveSpecification:
        AttributeName: "ttl"
        Enabled: true

  WriteCapacityScalableTarget:
    Type: AWS::ApplicationAutoScaling::ScalableTarget
    Properties:
      MaxCapacity: 15
      MinCapacity: 5
      ResourceId: !Join ["/", ["table", !Ref IdempotencyTable]]
      RoleARN: !GetAtt ScalingRole.Arn
      ScalableDimension: dynamodb:table:WriteCapacityUnits
      ServiceNamespace: dynamodb

  ScalingRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service:
                - application-autoscaling.amazonaws.com
            Action:
              - "sts:AssumeRole"
      Path: "/"
      ManagedPolicyArns:
        - !Ref IdempotencyDBAccessPolicy

  WriteScalingPolicy:
    Type: AWS::ApplicationAutoScaling::ScalingPolicy
    Properties:
      PolicyName: WriteAutoScalingPolicy
      PolicyType: TargetTrackingScaling
      ScalingTargetId: !Ref WriteCapacityScalableTarget
      TargetTrackingScalingPolicyConfiguration:
        TargetValue: 50.0
        ScaleInCooldown: 60
        ScaleOutCooldown: 60
        PredefinedMetricSpecification:
          PredefinedMetricType: DynamoDBWriteCapacityUtilization

  IdempotencyDBAccessPolicy:
    Type: AWS::IAM::ManagedPolicy
    Properties:
    PolicyDocument:
      Version: '2012-10-17'
      Statement:
        - Effect: Allow
          Action:
            - dynamodb:GetItem
            - dynamodb:PutItem
            - dynamodb:UpdateItem
          Resource: !GetAtt IdempotencyDBTable.Arn

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
        - "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"  # For CloudWatch Logs
        - "arn:aws:iam::aws:policy/CloudFrontFullAccess"  # For CloudFront access
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
                Resource: !GetAtt IdempotencyDBTable.Arn
        - PolicyName: SNSPublishAccess
          PolicyDocument:
            Version: '2012-10-17'
            Statement:
              - Effect: Allow
                Action:
                  - sns:Publish
                Resource: !Ref UserSubscribedTopic

  YourLambdaFunction:
    Type: AWS::Serverless::Function
    Properties:
      # ... other properties
      Role: !GetAtt SubscribeFunctionRole.Arn
```

## How It Works

1. When a request is received, the middleware generates an idempotency key from the request body.
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

## Error Handling

If the original handler returns an error, the middleware will:
1. Cache the error along with the response.
2. Return the same error for identical requests within the expiry period.

This ensures that error responses are also idempotent.

## TODO

- Add comprehensive benchmarks
- Add integration tests with DynamoDB or LocalStack
