package apigatewayresponse

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/scrambledeggs/booky-go-common/assert"
)

func TestSingleErrorResponseV2(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "*")
	status := http.StatusBadRequest
	origin := "http://localhost:3000"

	err := errors.New("invalid arguments")

	errorObj := ErrorResponseBody{
		Message: err.Error(),
		Code:    "INVALID_ARGUMENTS",
	}

	params := SingleErrorResponseV2Params{
		Status: status,
		Error:  errorObj,
		Origin: origin,
	}

	response, err := SingleErrorResponseV2(params)

	var responseBody ErrorResponseBody

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, responseBody, errorObj, "invalid value for error response")

	assert.Equal(t, response.StatusCode, status, "invalid status code")

	assert.Equal(t, response.Headers["Access-Control-Allow-Origin"], origin, "invalid origin")
}

func TestMultipleErrorsResponseV2(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "*")
	status := http.StatusInternalServerError
	origin := "http://localhost:3000"

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

	errors := []ErrorResponseBody{error1Obj, error2Obj}

	params := MultipleErrorResponseV2Params{
		Status: status,
		Errors: errors,
		Origin: origin,
	}

	response, err := MultipleErrorResponseV2(params)

	var responseBody MultipleErrorResponseBody

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, responseBody.Errors[0], error1Obj, "invalid error value for first element")
	assert.DeepEqual(t, responseBody.Errors[1], error2Obj, "invalid error value for second element")

	assert.Equal(t, response.StatusCode, status, "invalid status code")

	assert.Equal(t, response.Headers["Access-Control-Allow-Origin"], origin, "invalid origin")
}
