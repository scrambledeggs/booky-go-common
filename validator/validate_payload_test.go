package helpers

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestValidatePayload(t *testing.T) {
	cases := []struct {
		Name    string
		Request events.APIGatewayProxyRequest
		Schema  Schema
		WantErr bool
	}{
		{
			Name: "Valid query parameters",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"latitude":  "12.345678",
					"longitude": "123.456789",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"latitude":  {Type: Number, Required: true},
					"longitude": {Type: Number, Required: true},
				},
			},
			WantErr: false,
		},
		{
			Name: "Invalid number in query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"latitude":  "notanumber",
					"longitude": "123.456789",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"latitude":  {Type: Number, Required: true},
					"longitude": {Type: Number, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Missing required query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"latitude": "12.345678",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"latitude":  {Type: Number, Required: true},
					"longitude": {Type: Number, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Valid boolean in query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"active": "true",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"active": {Type: Boolean, Required: true},
				},
			},
			WantErr: false,
		},
		{
			Name: "Invalid boolean in query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"active": "notabool",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"active": {Type: Boolean, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Valid UUID in query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"id": "123e4567-e89b-12d3-a456-426614174000",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"id": {Type: UUID, Required: true},
				},
			},
			WantErr: false,
		},
		{
			Name: "Invalid UUID in query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"id": "notauuid",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"id": {Type: UUID, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Missing required query parameter",
			Request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"latitude": "12.345678",
				},
			},
			Schema: Schema{
				QueryParameters: &map[string]SchemaType{
					"latitude": {
						Type:     Number,
						Required: true,
					},
					"longitude": {
						Type:     Number,
						Required: true,
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "Valid body",
			Request: events.APIGatewayProxyRequest{
				Body: "{\"foo\":\"bar\"}",
			},
			Schema: Schema{
				Body: &map[string]SchemaType{
					"foo": {Type: String, Required: true},
				},
			},
			WantErr: false,
		},
		{
			Name: "Invalid JSON body",
			Request: events.APIGatewayProxyRequest{
				Body: "not a json",
			},
			Schema: Schema{
				Body: &map[string]SchemaType{
					"foo": {Type: String, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Missing required field in body",
			Request: events.APIGatewayProxyRequest{
				Body: "{}",
			},
			Schema: Schema{
				Body: &map[string]SchemaType{
					"foo": {Type: String, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Valid number in body",
			Request: events.APIGatewayProxyRequest{
				Body: "{\"num\":\"123.45\"}",
			},
			Schema: Schema{
				Body: &map[string]SchemaType{
					"num": {Type: Number, Required: true},
				},
			},
			WantErr: false,
		},
		{
			Name: "Invalid number in body",
			Request: events.APIGatewayProxyRequest{
				Body: "{\"num\":\"notanumber\"}",
			},
			Schema: Schema{
				Body: &map[string]SchemaType{
					"num": {Type: Number, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Valid path parameters",
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"id": "123",
				},
			},
			Schema: Schema{
				PathParameters: &map[string]SchemaType{
					"id": {Type: String, Required: true},
				},
			},
			WantErr: false,
		},
		{
			Name: "Missing required path parameter",
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{},
			},
			Schema: Schema{
				PathParameters: &map[string]SchemaType{
					"id": {Type: String, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Invalid UUID in path parameter",
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"id": "notauuid",
				},
			},
			Schema: Schema{
				PathParameters: &map[string]SchemaType{
					"id": {Type: UUID, Required: true},
				},
			},
			WantErr: true,
		},
		{
			Name: "Valid UUID in path parameter",
			Request: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{
					"id": "123e4567-e89b-12d3-a456-426614174000",
				},
			},
			Schema: Schema{
				PathParameters: &map[string]SchemaType{
					"id": {Type: UUID, Required: true},
				},
			},
			WantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := ValidatePayload(tc.Request, tc.Schema)
			if (err != nil) != tc.WantErr {
				t.Errorf("ValidatePayload() error = %v, wantErr %v", err, tc.WantErr)
			}
		})
	}
}
