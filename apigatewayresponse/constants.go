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
	// Read environment variables at runtime to support testing
	runtimeAllowOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	runtimeAllowMethods := os.Getenv("CORS_ALLOWED_METHODS")
	runtimeAllowHeaders := os.Getenv("CORS_ALLOWED_HEADERS")

	var headers = map[string]string{
		"Access-Control-Allow-Methods":     runtimeAllowMethods,
		"Access-Control-Allow-Headers":     runtimeAllowHeaders,
		"Access-Control-Allow-Credentials": "true",
		"Content-Type":                     "application/json",
	}

	if runtimeAllowOrigins == "*" {
		headers["Access-Control-Allow-Origin"] = origin
	} else if runtimeAllowOrigins != "" {
		allowed := strings.Split(runtimeAllowOrigins, ",")

		if slicesfunc.Contains(origin, allowed) {
			headers["Access-Control-Allow-Origin"] = origin
		}
	}

	return headers
}
