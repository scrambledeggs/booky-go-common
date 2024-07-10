package apigatewayresponse

import (
	"encoding/json"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
)

func SuccessResponse(status int, data ...any) (events.APIGatewayProxyResponse, error) {
	var metadata any = nil
	var strBody []byte

	response := events.APIGatewayProxyResponse{
		Headers:    HttpHeaders,
		StatusCode: status,
	}

	isSlice := IsSlice(data[0])

	if isSlice {
		body := SuccessResponseBody{
			Results: data[0],
		}

		if len(data) > 1 {
			metadata = data[1]
		}

		body.Metadata = metadata

		strBody, _ = json.Marshal(body)

	} else {
		strBody, _ = json.Marshal(data[0])
	}

	response.Body = string(strBody)

	return response, nil
}

func IsSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
