package apigatewayresponse

import (
	"encoding/json"
	"math"

	"github.com/aws/aws-lambda-go/events"
)

type SingleSuccessResponseV2Params struct {
	Origin            string
	Status            int
	Data              any
	Headers           map[string]string
	MultiValueHeaders map[string][]string
}

func SingleSuccessResponseV2(params SingleSuccessResponseV2Params) (events.APIGatewayProxyResponse, error) {
	var strBody []byte

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

	strBody, err := json.Marshal(params.Data)

	if err != nil {
		panic(err.Error())
	}

	response.Body = string(strBody)

	return response, nil
}

type MultipleSuccessResponseV2Params struct {
	Origin            string
	Status            int
	Data              any
	Metadata          any
	Headers           map[string]string
	MultiValueHeaders map[string][]string
}

func MultipleSuccessResponseV2(params MultipleSuccessResponseV2Params) (events.APIGatewayProxyResponse, error) {
	var strBody []byte

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

	body := MultipleSuccessResponseBody{
		Results: params.Data,
	}

	if params.Metadata != nil {
		mtdt, _ := params.Metadata.(map[string]any)

		totalCount, isTotalCountTypeInt64 := mtdt["total_count"].(int64)
		resultsPerPage, isResultsPerPageTypeInt64 := mtdt["results_per_page"].(int64)

		if isResultsPerPageTypeInt64 && isTotalCountTypeInt64 {
			maxPage := int64(math.Ceil(float64(totalCount) / float64(resultsPerPage)))

			mtdt["max_page"] = maxPage
		}

		body.Metadata = mtdt
	}

	strBody, err := json.Marshal(body)

	if err != nil {
		panic(err.Error())
	}

	response.Body = string(strBody)

	return response, nil
}
