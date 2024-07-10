package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type MultipleSuccessResponseBody struct {
	Results  any `json:"results"`
	Metadata any `json:"metadata"`
}

func SingleSuccessResponse(status int, data any) (events.APIGatewayProxyResponse, error) {
	var strBody []byte

	response := events.APIGatewayProxyResponse{
		Headers:    HTTPHeaders,
		StatusCode: status,
	}

	strBody, err := json.Marshal(data)

	if err != nil {
		panic(err.Error())
	}

	response.Body = string(strBody)

	return response, nil
}

func MultipleSuccessResponse(status int, data any, metadata any) (events.APIGatewayProxyResponse, error) {
	var strBody []byte

	response := events.APIGatewayProxyResponse{
		Headers:    HTTPHeaders,
		StatusCode: status,
	}

	body := MultipleSuccessResponseBody{
		Results: data,
	}

	if metadata != nil {
		body.Metadata = metadata
	}

	strBody, err := json.Marshal(body)

	if err != nil {
		panic(err.Error())
	}

	response.Body = string(strBody)

	return response, nil
}
