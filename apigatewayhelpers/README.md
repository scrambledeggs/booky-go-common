# API Gateway Helpers

## Initialize pagination info with metadata

Pass `request.QueryStringParameters` from lambda request

We can use the function as follows
```go
offset, resultsPerPage, metadata := apigatewayhelpers.InitPaginate(request.QueryStringParameters)
```
