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

    // Wrap your handler with the idempotency middleware and custom options
    lambda.Start(idempotency.NewIdempotentHandlerWithOptions(handler, options))
}
```

## Configuration

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
  - `status` (String): Request status (in_progress, completed, expired)
  - `expiration` (String): Expiration timestamp
  - `error` (String): Error message if applicable
  - `request_headers` (String): Original request headers

### CloudFormation Template

You can use the following CloudFormation template in your `template.yaml` file to create the required DynamoDB table with auto-scaling capabilities:

```yaml
Resources:
  IdempotencyDBTable:
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
        AttributeName: "expiration"
        Enabled: true

  WriteCapacityScalableTarget:
    Type: AWS::ApplicationAutoScaling::ScalableTarget
    Properties:
      MaxCapacity: 15
      MinCapacity: 5
      ResourceId: !Join ["/", ["table", !Ref IdempotencyDBTable]]
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
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Action:
              - "dynamodb:DescribeTable"
              - "dynamodb:UpdateTable"
              - "dynamodb:GetItem"
              - "dynamodb:PutItem"
              - "dynamodb:DeleteItem"
              - "dynamodb:UpdateItem"
            Resource: !GetAtt IdempotencyDBTable.Arn

  YourLambdaFunction:
    Type: AWS::Serverless::Function
    Properties:
      # ... other properties
      Role: !GetAtt LambdaExecutionRole.Arn
      Environment:
        Variables:
          IDEMPOTENCY_DB_TABLE: !Ref IdempotencyDBTable
          AMZ_REGION: !Ref AWS::Region

  LambdaExecutionRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
          - Effect: Allow
            Principal:
              Service: lambda.amazonaws.com
            Action: sts:AssumeRole
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
        - !Ref IdempotencyDBAccessPolicy
```

This template includes:

1. **Auto-scaling for DynamoDB**: Automatically scales write capacity between 5-15 units based on a 50% utilization target
2. **IAM Roles and Policies**: Proper permissions for both the scaling service and your Lambda function
3. **Scaling Configuration**: Cooldown periods to prevent rapid scaling changes

For additional production optimizations, consider adding:

```yaml
  IdempotencyDBTable:
    Type: AWS::DynamoDB::Table
    Properties:
      # ... existing properties
      PointInTimeRecoverySpecification:
        PointInTimeRecoveryEnabled: true
      BillingMode: PROVISIONED  # Or PAY_PER_REQUEST for on-demand capacity
```

## Environment Variables

The following environment variables are used by the idempotency package:

- `IDEMPOTENCY_DB_TABLE`: DynamoDB table name (required)
- `AMZ_REGION`: AWS region (required)

> **Note:** The idempotency package uses IAM roles for AWS authentication. Make sure your Lambda function has the appropriate IAM role with permissions to access DynamoDB as shown in the CloudFormation template example.

## How It Works

1. When a request is received, the middleware calculates an idempotency key based on the request body
2. It checks if a record with the same key exists in DynamoDB
3. If a record exists and is completed, it returns the stored response
4. If a record exists and is in progress, it returns a 425 Too Early status
5. If no record exists or the record is expired, it creates a new record, processes the request, and stores the response

## Error Handling

The middleware preserves both successful responses and errors. If the original request resulted in an error, the same error will be returned for duplicate requests.
