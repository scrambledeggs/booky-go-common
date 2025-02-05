package authctxparser

import (
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type RequestContextUser struct {
	ID             *pgtype.UUID `json:"user_id"`
	MobileNumber   *string      `json:"user_mobile_number"`
	Name           *string      `json:"user_name"`
	Email          *string      `json:"user_email"`
	MayaCustomerID *string      `json:"user_maya_customer_id"`
	DeactivatedAt  *string      `json:"user_deactivated_at"`
}

var (
	ErrInvalidContext = errors.New("invalid context")
	ErrInvalidUser    = errors.New("invalid user")
)

func ParseUser(rawAuthCtx map[string]interface{}) (*RequestContextUser, error) {
	authCtxStr, err := json.Marshal(rawAuthCtx)

	if err != nil {
		return nil, ErrInvalidContext
	}

	var user *RequestContextUser

	if err := json.Unmarshal(authCtxStr, &user); err != nil {
		return nil, ErrInvalidUser
	}

	if user == nil || user.ID == nil {
		return nil, ErrInvalidUser
	}

	return user, nil
}
