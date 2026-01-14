package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type MultipleErrorResponseV2Params struct {
	Origin            string
	Status            int
	Errors            []ErrorResponseBody
	Headers           map[string]string
	MultiValueHeaders map[string][]string
}

func MultipleErrorResponseV2(params MultipleErrorResponseV2Params) (events.APIGatewayProxyResponse, error) {
	headers := buildResponseHeaders(params.Origin)

	if params.Headers != nil {
		for k, v := range params.Headers {
			headers[k] = v
		}
	}

	response := events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: params.Status,
	}

	if params.MultiValueHeaders != nil {
		response.MultiValueHeaders = params.MultiValueHeaders
	}

	body := MultipleErrorResponseBody{
		Errors: params.Errors,
	}

	strBody, err := json.Marshal(body)

	if err != nil {
		panic(err.Error())
	}

	response.Body = string(strBody)

	return response, nil
}

type SingleErrorResponseV2Params struct {
	Origin            string
	Status            int
	Error             ErrorResponseBody
	Headers           map[string]string
	MultiValueHeaders map[string][]string
}

func SingleErrorResponseV2(params SingleErrorResponseV2Params) (events.APIGatewayProxyResponse, error) {
	headers := buildResponseHeaders(params.Origin)

	if params.Headers != nil {
		for k, v := range params.Headers {
			headers[k] = v
		}
	}

	response := events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: params.Status,
	}

	if params.MultiValueHeaders != nil {
		response.MultiValueHeaders = params.MultiValueHeaders
	}

	strBody, er := json.Marshal(params.Error)

	if er != nil {
		panic(er.Error())
	}

	response.Body = string(strBody)

	return response, nil
}
