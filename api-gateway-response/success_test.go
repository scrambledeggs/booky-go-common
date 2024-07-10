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
		PageCount:   10,
		ResultCount: 100,
	}

	response, err := MultipleSuccessResponse(status, body, metadata)

	var responseBody SuccessResponse
	json.Unmarshal([]byte(response.Body), &responseBody)

	metadataStr, _ := json.Marshal(responseBody.Metadata)
	var metadataRes PaginationMetadata
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, _ := json.Marshal(responseBody.Results)
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
		PageCount:   10,
		ResultCount: 100,
	}

	multiResponse, _ := MultipleSuccessResponse(status, multipleBody, metadata)

	var responseBody SuccessResponse
	json.Unmarshal([]byte(multiResponse.Body), &responseBody)

	metadataStr, _ := json.Marshal(responseBody.Metadata)
	var metadataRes PaginationMetadata
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, _ := json.Marshal(responseBody.Results)
	var resultsRes []map[string]string
	json.Unmarshal([]byte(resultsStr), &resultsRes)

	fmt.Println(multiResponse.Body, multiResponse.StatusCode)

	// Output:
	// {"results":[{"naknang":"sonof","patatas":"potato"},{"potato":"patatas","sonof":"naknang"}],"metadata":{"page_count":10,"result_count":100}} 200
}
