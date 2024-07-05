package apigatewayresponse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponseSingleResult(t *testing.T) {
	status := http.StatusOK

	body := map[string]string{
		"naknang": "sonof",
		"patatas": "potato",
	}

	response, err := SuccessResponse(status, body)

	var responseBody map[string]string

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.Equal(t, nil, err)
	assert.Equal(t, responseBody["naknang"], body["naknang"])
	assert.Equal(t, responseBody["patatas"], body["patatas"])
	assert.Equal(t, response.StatusCode, status)
}

func TestSuccessResponseMultipleResults(t *testing.T) {
	status := http.StatusOK

	body := []map[string]string{
		{"naknang": "sonof", "patatas": "potato"},
		{"sonof": "naknang", "potato": "patatas"},
	}

	metadata := PaginationMetadata{
		MaxPage:        10,
		Page:           1,
		ResultsPerPage: 10,
		TotalCount:     100,
	}

	response, err := SuccessResponse(status, body, metadata)

	var responseBody SuccessResponseBody
	json.Unmarshal([]byte(response.Body), &responseBody)

	metadataStr, _ := json.Marshal(responseBody.Metadata)
	var metadataRes PaginationMetadata
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, _ := json.Marshal(responseBody.Results)
	var resultsRes []map[string]string
	json.Unmarshal([]byte(resultsStr), &resultsRes)

	assert.Equal(t, nil, err)

	assert.Equal(t, metadataRes.MaxPage, metadata.MaxPage)
	assert.Equal(t, metadataRes.Page, metadata.Page)
	assert.Equal(t, metadataRes.ResultsPerPage, metadata.ResultsPerPage)
	assert.Equal(t, metadataRes.TotalCount, metadata.TotalCount)

	assert.Equal(t, resultsRes[0]["naknang"], body[0]["naknang"])
	assert.Equal(t, resultsRes[0]["patatas"], body[0]["patatas"])
	assert.Equal(t, resultsRes[1]["sonof"], body[1]["sonof"])
	assert.Equal(t, resultsRes[1]["potato"], body[1]["potato"])
	assert.Equal(t, response.StatusCode, status)
}

func ExampleSuccessResponse() {
	status := http.StatusOK

	singleBody := map[string]string{
		"naknang": "sonof",
		"patatas": "potato",
	}

	singleResponse, _ := SuccessResponse(status, singleBody)

	multipleBody := []map[string]string{
		{"naknang": "sonof", "patatas": "potato"},
		{"sonof": "naknang", "potato": "patatas"},
	}

	metadata := PaginationMetadata{
		MaxPage:        10,
		Page:           1,
		ResultsPerPage: 10,
		TotalCount:     100,
	}

	multiResponse, _ := SuccessResponse(status, multipleBody, metadata)

	var responseBody SuccessResponseBody
	json.Unmarshal([]byte(multiResponse.Body), &responseBody)

	metadataStr, _ := json.Marshal(responseBody.Metadata)
	var metadataRes PaginationMetadata
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, _ := json.Marshal(responseBody.Results)
	var resultsRes []map[string]string
	json.Unmarshal([]byte(resultsStr), &resultsRes)

	fmt.Println(singleResponse.Body, singleResponse.StatusCode)
	fmt.Println(multiResponse.Body, multiResponse.StatusCode)

	// Output:
	// {"naknang":"sonof","patatas":"potato"} 200
	// {"results":[{"naknang":"sonof","patatas":"potato"},{"potato":"patatas","sonof":"naknang"}],"metadata":{"max_page":10,"page":1,"results_per_page":10,"total_count":100}} 200
}
