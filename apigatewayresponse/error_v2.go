package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func MultipleErrorResponseV2(origin string, status int, errors []ErrorResponseBody) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers:    buildResponseHeaders(origin),
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

func SingleErrorResponseV2(origin string, status int, err ErrorResponseBody) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers:    buildResponseHeaders(origin),
		StatusCode: status,
	}

	strBody, er := json.Marshal(err)

	if er != nil {
		panic(er.Error())
	}

	response.Body = string(strBody)

	return response, nil
}
