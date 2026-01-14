package apigatewayresponse

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type MultipleErrorResponseV2Params struct {
	Origin  string
	Status  int
	Errors  []ErrorResponseBody
	Headers map[string]string
	Cookies []string
}

func MultipleErrorResponseV2(params MultipleErrorResponseV2Params) (events.APIGatewayV2HTTPResponse, error) {
	headers := buildResponseHeaders(params.Origin)

	if params.Headers != nil {
		for k, v := range params.Headers {
			headers[k] = v
		}
	}

	response := events.APIGatewayV2HTTPResponse{
		Headers:    headers,
		StatusCode: params.Status,
	}

	if params.Cookies != nil {
		response.Cookies = params.Cookies
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
	Origin  string
	Status  int
	Error   ErrorResponseBody
	Headers map[string]string
	Cookies []string
}

func SingleErrorResponseV2(params SingleErrorResponseV2Params) (events.APIGatewayV2HTTPResponse, error) {
	headers := buildResponseHeaders(params.Origin)

	if params.Headers != nil {
		for k, v := range params.Headers {
			headers[k] = v
		}
	}

	response := events.APIGatewayV2HTTPResponse{
		Headers:    headers,
		StatusCode: params.Status,
	}

	if params.Cookies != nil {
		response.Cookies = params.Cookies
	}

	strBody, er := json.Marshal(params.Error)

	if er != nil {
		panic(er.Error())
	}

	response.Body = string(strBody)

	return response, nil
}
