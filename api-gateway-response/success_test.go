package apigatewayresponse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/scrambledeggs/booky-go-common/assert"
)

func TestSingleSuccessResponse(t *testing.T) {
	status := http.StatusOK

	body := map[string]string{
		"naknang": "sonof",
		"patatas": "potato",
	}

	response, err := SingleSuccessResponse(status, body)

	var responseBody map[string]string

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, responseBody, body, "invalid value for body")

	assert.Equal(t, response.StatusCode, status, "invalid status code")
}

func TestMultipleSuccessResponse(t *testing.T) {
	status := http.StatusOK

	body := []map[string]string{
		{"naknang": "sonof", "patatas": "potato"},
		{"sonof": "naknang", "potato": "patatas"},
	}

	metadata := PaginationMetadata{
		Page:           1,
		ResultsPerPage: 10,
		MaxPage:        10,
		TotalCount:     100,
	}

	response, err := MultipleSuccessResponse(status, body, metadata)

	if err != nil {
		panic(err.Error())
	}

	var responseBody MultipleSuccessResponseBody
	json.Unmarshal([]byte(response.Body), &responseBody)

	metadataStr, err := json.Marshal(responseBody.Metadata)

	if err != nil {
		panic(err.Error())
	}

	var metadataRes PaginationMetadata
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, err := json.Marshal(responseBody.Results)

	if err != nil {
		panic(err.Error())
	}

	var resultsRes []map[string]string
	json.Unmarshal([]byte(resultsStr), &resultsRes)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, metadataRes, metadata, "invalid metadata")
	assert.DeepEqual(t, resultsRes[0], body[0], "invalid value first element")
	assert.DeepEqual(t, resultsRes[1], body[1], "invalid value for second element")

	assert.Equal(t, response.StatusCode, status, "invalid status code")
}

func ExampleSingleSuccessResponse() {
	status := http.StatusOK

	singleBody := map[string]string{
		"naknang": "sonof",
		"patatas": "potato",
	}

	singleResponse, _ := SingleSuccessResponse(status, singleBody)

	fmt.Println(singleResponse.Body, singleResponse.StatusCode)

	// Output:
	// {"naknang":"sonof","patatas":"potato"} 200
}

func ExampleMultipleSuccessResponse() {
	status := http.StatusOK

	multipleBody := []map[string]string{
		{"naknang": "sonof", "patatas": "potato"},
		{"sonof": "naknang", "potato": "patatas"},
	}

	metadata := PaginationMetadata{
		Page:           1,
		ResultsPerPage: 10,
		MaxPage:        10,
		TotalCount:     100,
	}

	multiResponse, _ := MultipleSuccessResponse(status, multipleBody, metadata)

	var responseBody MultipleSuccessResponseBody
	json.Unmarshal([]byte(multiResponse.Body), &responseBody)

	metadataStr, err := json.Marshal(responseBody.Metadata)

	if err != nil {
		panic(err.Error())
	}

	var metadataRes PaginationMetadata
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, err := json.Marshal(responseBody.Results)

	if err != nil {
		panic(err.Error())
	}

	var resultsRes []map[string]string
	json.Unmarshal([]byte(resultsStr), &resultsRes)

	fmt.Println(multiResponse.Body, multiResponse.StatusCode)

	// Output:
	// {"results":[{"naknang":"sonof","patatas":"potato"},{"potato":"patatas","sonof":"naknang"}],"metadata":{"page":1,"results_per_page":10,"total_count":100,"max_page":10}} 200
}
