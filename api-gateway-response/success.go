package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func SingleSuccessResponse(status int, data any) (events.APIGatewayProxyResponse, error) {
	var strBody []byte

	response := events.APIGatewayProxyResponse{
		Headers:    HttpHeaders,
		StatusCode: status,
	}

	strBody, _ = json.Marshal(data)

	response.Body = string(strBody)

	return response, nil
}

func MultipleSuccessResponse(status int, data any, metadata any) (events.APIGatewayProxyResponse, error) {
	var strBody []byte

	response := events.APIGatewayProxyResponse{
		Headers:    HttpHeaders,
		StatusCode: status,
	}

	body := SuccessResponse{
		Results: data,
	}

	if metadata != nil {
		body.Metadata = metadata
	}

	strBody, _ = json.Marshal(body)

	response.Body = string(strBody)

	return response, nil
}
