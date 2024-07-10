package apigatewayresponse

import "os"

var allowHeaders = os.Getenv("CORS_ALLOWED_HEADERS")
var allowOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")
var allowMethods = os.Getenv("CORS_ALLOWED_METHODS")

var PreflightHttpHeaders = map[string]string{
	"Access-Control-Allow-Origin":  allowOrigins,
	"Access-Control-Allow-Methods": "OPTIONS",
	"Access-Control-Allow-Headers": allowHeaders,
}

var HttpHeaders = map[string]string{
	"Access-Control-Allow-Origin":  allowOrigins,
	"Access-Control-Allow-Methods": allowMethods,
	"Access-Control-Allow-Headers": allowHeaders,
	"Content-Type":                 "application/json",
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SuccessResponseBody struct {
	Results  any `json:"results"`
	Metadata any `json:"metadata"`
}

type PaginationMetadata struct {
	PageCount   int `json:"max_page"`
	ResultCount int `json:"results_per_page"`
}
