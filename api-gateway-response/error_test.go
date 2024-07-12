package apigatewayresponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/scrambledeggs/booky-go-common/assert"
)

func TestSingleErrorResponse(t *testing.T) {
	status := http.StatusBadRequest

	err := errors.New("invalid arguments")

	errorObj := ErrorResponseBody{
		Message: err.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, err := SingleErrorResponse(status, errorObj)

	var responseBody ErrorResponseBody

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, responseBody, errorObj, "invalid value for error response")

	assert.Equal(t, response.StatusCode, status, "invalid status code")
}

func TestMultipleErrorsResponse(t *testing.T) {
	status := http.StatusInternalServerError

	err1 := errors.New("invalid name")

	error1Obj := ErrorResponseBody{
		Message: err1.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	err2 := errors.New("invalid slug")

	error2Obj := ErrorResponseBody{
		Message: err2.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, err := MultipleErrorResponse(status, []ErrorResponseBody{error1Obj, error2Obj})

	var responseBody MultipleErrorResponseBody

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, responseBody.Errors[0], error1Obj, "invalid error value for first element")
	assert.DeepEqual(t, responseBody.Errors[1], error2Obj, "invalid error value for second element")

	assert.Equal(t, response.StatusCode, status, "invalid status code")
}

func ExampleSingleErrorResponse() {
	status := http.StatusBadRequest

	err := errors.New("invalid arguments")

	errorObj := ErrorResponseBody{
		Message: err.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, _ := SingleErrorResponse(status, errorObj)

	fmt.Println(response.Body, response.StatusCode)

	// Output: {"message":"invalid arguments","code":"INVALID_ARGUMENTS"} 400
}

func ExampleMultipleErrorsResponse() {
	status := http.StatusInternalServerError

	err1 := errors.New("invalid name")

	error1Obj := ErrorResponseBody{
		Message: err1.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	err2 := errors.New("invalid slug")

	error2Obj := ErrorResponseBody{
		Message: err2.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, _ := MultipleErrorResponse(status, []ErrorResponseBody{error1Obj, error2Obj})

	fmt.Println(response.Body, response.StatusCode)

	// Output: {"errors":[{"message":"invalid name","code":"INVALID_ARGUMENTS"},{"message":"invalid slug","code":"INVALID_ARGUMENTS"}]} 500
}
