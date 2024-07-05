package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type MultipleErrorsResponseType struct {
	Errors []ErrorResponse `json:"errors"`
}

func MultipleErrorsResponse(status int, errors []ErrorResponse) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers:    HttpHeaders,
		StatusCode: status,
	}

	body := MultipleErrorsResponseType{
		Errors: errors,
	}

	rData, _ := json.Marshal(body)

	response.Body = string(rData)

	return response, nil
}

func SingleErrorResponse(status int, err ErrorResponse) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers:    HttpHeaders,
		StatusCode: status,
	}

	rData, _ := json.Marshal(err)

	response.Body = string(rData)

	return response, nil
}
