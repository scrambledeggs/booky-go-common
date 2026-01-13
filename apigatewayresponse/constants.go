package apigatewayresponse

import (
	"fmt"
	"os"
	"strings"

	"github.com/scrambledeggs/booky-go-common/slicesfunc"
)

var HTTPHeaders = map[string]string{
	"Access-Control-Allow-Origin":  os.Getenv("CORS_ALLOWED_ORIGINS"),
	"Access-Control-Allow-Methods": os.Getenv("CORS_ALLOWED_METHODS"),
	"Access-Control-Allow-Headers": os.Getenv("CORS_ALLOWED_HEADERS"),
	"Content-Type":                 "application/json",
}

type PaginationMetadata struct {
	Page           int `json:"page"`
	ResultsPerPage int `json:"results_per_page"`
	TotalCount     int `json:"total_count"`
}

func buildResponseHeaders(origin string) map[string]string {
	var headers = map[string]string{
		"Access-Control-Allow-Methods": os.Getenv("CORS_ALLOWED_METHODS"),
		"Access-Control-Allow-Headers": os.Getenv("CORS_ALLOWED_HEADERS"),
		"Content-Type":                 "application/json",
	}

	allowOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	fmt.Println("===============================")
	fmt.Println("allowOrigins", allowOrigins)
	fmt.Println("origin", origin)

	if allowOrigins == "*" {
		headers["Access-Control-Allow-Origin"] = origin // ❗ NOT "*"
	} else {
		allowed := strings.Split(allowOrigins, ",")
		if slicesfunc.Contains(origin, allowed) {
			headers["Access-Control-Allow-Origin"] = origin
		}
	}

	fmt.Println("responseHeaders", headers)

	return headers
}
