# API Gateway Helpers

## Pagination Info Parser

Pass `request.RequestContext.Authorizer` from lambda request. This requires the header `x-user-token` to be included in the request.

We can use the function as follows
```go
  user, err := authctxparser.ParseUser(request.RequestContext.Authorizer)
```

`user` variable is an instance of

```go
type RequestContextUser struct {
	ID             *pgtype.UUID `json:"user_id"`
	MobileNumber   *string      `json:"user_mobile_number"`
	Name           *string      `json:"user_name"`
	Email          *string      `json:"user_email"`
	MayaCustomerID *string      `json:"user_maya_customer_id"`
	DeactivatedAt  *string      `json:"user_deactivated_at"`
}
```
