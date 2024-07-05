package apigatewayresponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleErrorResponse(t *testing.T) {
	status := http.StatusBadRequest

	err := errors.New("invalid arguments")

	errorObj := ErrorResponse{
		Message: err.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, err := SingleErrorResponse(status, errorObj)

	var responseBody ErrorResponse

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.Equal(t, err, nil)

	assert.Equal(t, responseBody.Code, errorObj.Code)
	assert.Equal(t, responseBody.Message, errorObj.Message)
}

func TestMultipleErrorsResponse(t *testing.T) {
	status := http.StatusInternalServerError

	err1 := errors.New("invalid name")

	error1Obj := ErrorResponse{
		Message: err1.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	err2 := errors.New("invalid slug")

	error2Obj := ErrorResponse{
		Message: err2.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, err := MultipleErrorsResponse(status, []ErrorResponse{error1Obj, error2Obj})

	var responseBody MultipleErrorsResponseType

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.Equal(t, err, nil)

	assert.Equal(t, responseBody.Errors[0].Code, error1Obj.Code)
	assert.Equal(t, responseBody.Errors[0].Message, error1Obj.Message)
	assert.Equal(t, responseBody.Errors[1].Code, error2Obj.Code)
	assert.Equal(t, responseBody.Errors[1].Message, error2Obj.Message)
}

func ExampleSingleErrorResponse() {
	status := http.StatusBadRequest

	err := errors.New("invalid arguments")

	errorObj := ErrorResponse{
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

	error1Obj := ErrorResponse{
		Message: err1.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	err2 := errors.New("invalid slug")

	error2Obj := ErrorResponse{
		Message: err2.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	response, _ := MultipleErrorsResponse(status, []ErrorResponse{error1Obj, error2Obj})

	fmt.Println(response.Body, response.StatusCode)

	// Output: {"errors":[{"message":"invalid name","code":"INVALID_ARGUMENTS"},{"message":"invalid slug","code":"INVALID_ARGUMENTS"}]} 500
}
