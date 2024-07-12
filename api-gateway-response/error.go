package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponseBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type MultipleErrorResponseBody struct {
	Errors []ErrorResponseBody `json:"errors"`
}

func MultipleErrorResponse(status int, errors []ErrorResponseBody) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers:    HTTPHeaders,
		StatusCode: status,
	}

	body := MultipleErrorResponseBody{
		Errors: errors,
	}

	strBody, err := json.Marshal(body)

	if err != nil {
		panic(err.Error())
	}

	response.Body = string(strBody)

	return response, nil
}

func SingleErrorResponse(status int, err ErrorResponseBody) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers:    HTTPHeaders,
		StatusCode: status,
	}

	strBody, er := json.Marshal(err)

	if er != nil {
		panic(er.Error())
	}

	response.Body = string(strBody)

	return response, nil
}
