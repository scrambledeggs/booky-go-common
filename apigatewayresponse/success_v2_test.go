package apigatewayresponse

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/scrambledeggs/booky-go-common/assert"
)

func TestSingleSuccessResponseV2(t *testing.T) {
	status := http.StatusOK
	origin := "http://localhost:3000"

	body := map[string]string{
		"naknang": "sonof",
		"patatas": "potato",
	}

	params := SingleSuccessResponseV2Params{
		Status: status,
		Data:   body,
		Origin: origin,
	}

	response, err := SingleSuccessResponseV2(params)

	var responseBody map[string]string

	json.Unmarshal([]byte(response.Body), &responseBody)

	assert.ShouldBeNil(t, err)

	assert.DeepEqual(t, responseBody, body, "invalid value for body")

	assert.Equal(t, response.StatusCode, status, "invalid status code")

	assert.Equal(t, response.Headers["Access-Control-Allow-Origin"], origin, "invalid origin")
}

func TestMultipleSuccessResponseV2(t *testing.T) {
	status := http.StatusOK
	origin := "http://localhost:3000"

	body := []map[string]string{
		{"naknang": "sonof", "patatas": "potato"},
		{"sonof": "naknang", "potato": "patatas"},
	}

	metadata := map[string]any{
		"page":             int64(1),
		"results_per_page": int64(10),
		"total_count":      int64(100),
	}

	params := MultipleSuccessResponseV2Params{
		Status:   status,
		Data:     body,
		Metadata: metadata,
		Origin:   origin,
	}

	response, err := MultipleSuccessResponseV2(params)

	if err != nil {
		panic(err.Error())
	}

	var responseBody MultipleSuccessResponseBody
	json.Unmarshal([]byte(response.Body), &responseBody)

	metadataStr, err := json.Marshal(responseBody.Metadata)

	if err != nil {
		panic(err.Error())
	}

	var metadataRes map[string]int64
	json.Unmarshal([]byte(metadataStr), &metadataRes)

	resultsStr, err := json.Marshal(responseBody.Results)

	if err != nil {
		panic(err.Error())
	}

	var resultsRes []map[string]string
	json.Unmarshal([]byte(resultsStr), &resultsRes)

	assert.ShouldBeNil(t, err)

	assert.Equal(t, metadataRes["page"], metadata["page"], "invalid metadata")
	assert.Equal(t, metadataRes["results_per_page"], metadata["results_per_page"], "invalid metadata")
	assert.Equal(t, metadataRes["total_count"], metadata["total_count"], "invalid metadata")
	assert.Equal(t, metadataRes["max_page"], int64(10), "invalid metadata")
	assert.DeepEqual(t, resultsRes[0], body[0], "invalid value first element")
	assert.DeepEqual(t, resultsRes[1], body[1], "invalid value for second element")

	assert.Equal(t, response.StatusCode, status, "invalid status code")

	assert.Equal(t, response.Headers["Access-Control-Allow-Origin"], origin, "invalid origin")
}
