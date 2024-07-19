package apigatewayresponse

import (
	"encoding/json"
	"math"

	"github.com/aws/aws-lambda-go/events"
)

type MultipleSuccessResponseBody struct {
	Results  any `json:"results"`
	Metadata any `json:"metadata,omitempty"`
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
		mtdt, _ := metadata.(map[string]any)

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
