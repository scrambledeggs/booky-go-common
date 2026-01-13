package apigatewayresponse

import (
	"os"
	"strings"

	"github.com/scrambledeggs/booky-go-common/slicesfunc"
)

var allowHeaders = os.Getenv("CORS_ALLOWED_HEADERS")
var allowOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")
var allowMethods = os.Getenv("CORS_ALLOWED_METHODS")

var HTTPHeaders = map[string]string{
	"Access-Control-Allow-Origin":  allowOrigins,
	"Access-Control-Allow-Methods": allowMethods,
	"Access-Control-Allow-Headers": allowHeaders,
	"Content-Type":                 "application/json",
}

type PaginationMetadata struct {
	Page           int `json:"page"`
	ResultsPerPage int `json:"results_per_page"`
	TotalCount     int `json:"total_count"`
}

func buildResponseHeaders(origin string) map[string]string {
	var headers = map[string]string{
		"Access-Control-Allow-Methods":     allowMethods,
		"Access-Control-Allow-Headers":     allowHeaders,
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}

	if allowOrigins == "*" {
		headers["Access-Control-Allow-Origin"] = origin
	} else {
		allowed := strings.Split(allowOrigins, ",")
		if slicesfunc.Contains(origin, allowed) {
			headers["Access-Control-Allow-Origin"] = origin
		}
	}

	return headers
}
