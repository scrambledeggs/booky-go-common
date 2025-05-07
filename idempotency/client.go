package idempotency

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type idempotencyDBClient struct {
	*dynamodb.Client
	tableName *string
	context   context.Context
}

func newIdempotencyDBClient(ctx context.Context, tableName string) idempotencyDBClient {
	client := NewDynamoDBClient(ctx, nil)
	return idempotencyDBClient{Client: client, tableName: &tableName, context: ctx}
}

func newIdempotencyDBClientWithUrl(ctx context.Context, tableName, dbUrl string) idempotencyDBClient {
	client := NewDynamoDBClient(ctx, &dbUrl)
	return idempotencyDBClient{Client: client, tableName: &tableName, context: ctx}
}

func (idbc idempotencyDBClient) Put(record idempotencyRecord) {
	_, err := idbc.PutItem(idbc.context, &dynamodb.PutItemInput{Item: record.Item(), TableName: idbc.tableName})
	if err != nil {
		panic(err.Error())
	}
}

func (idbc idempotencyDBClient) Get(query idempotencyRecord) *idempotencyRecord {
	output, err := idbc.GetItem(idbc.context, &dynamodb.GetItemInput{Key: query.GetKey(), TableName: idbc.tableName})

	if err != nil {
		panic(err.Error())
	}

	if output.Item == nil {
		return nil
	}

	var result *idempotencyRecord
	err = attributevalue.UnmarshalMap(output.Item, &result)

	if err != nil {
		panic(err.Error())
	}

	return result
}
