package authctxparser

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scrambledeggs/booky-go-common/converters"
)

func TestParseUser(t *testing.T) {
	rawBlankAuthorizerContext := map[string]interface{}{}

	rawAuthorizerContext := map[string]interface{}{
		"user_email":            "naknang@patatas.com",
		"user_id":               "0496a7eb-d7a9-4ecc-9276-1d1af4dc9158",
		"user_name":             "Naknang Patatas",
		"user_mobile_number":    "09999999999",
		"user_birth_date":       "1991-01-01",
		"user_gender":           "prefer not to say",
		"user_maya_customer_id": nil,
		"user_deactivated_at":   nil,
	}

	userMobileNumber := rawAuthorizerContext["user_mobile_number"].(string)
	userName := rawAuthorizerContext["user_name"].(string)
	userEmail := rawAuthorizerContext["user_email"].(string)
	userBirthDate := rawAuthorizerContext["user_birth_date"].(string)
	userGender := rawAuthorizerContext["user_gender"].(string)

	rawExpectedUser := RequestContextUser{
		ID:             &pgtype.UUID{Bytes: uuid.MustParse(rawAuthorizerContext["user_id"].(string)), Valid: true},
		MobileNumber:   &userMobileNumber,
		Name:           &userName,
		Email:          &userEmail,
		BirthDate:      &userBirthDate,
		Gender:         &userGender,
		MayaCustomerID: nil,
		DeactivatedAt:  nil,
	}

	expectedUser := converters.StructToInterface(rawExpectedUser)

	t.Run("should return an error", func(t *testing.T) {
		rawGot, err := ParseUser(rawBlankAuthorizerContext)

		if rawGot != nil {
			t.Errorf("Expected user %v, but got %v", nil, rawGot)
		}

		if err != ErrInvalidUser {
			t.Errorf("Expected message %s, but got %s", ErrInvalidUser, err)
		}
	})

	t.Run("should return the user from authorizer context", func(t *testing.T) {
		rawGot, errorResp := ParseUser(rawAuthorizerContext)

		got := converters.StructToInterface(rawGot)

		if !reflect.DeepEqual(expectedUser, got) {
			t.Errorf("Expected user %v, but got %v", expectedUser, got)
		}

		if errorResp != nil {
			t.Errorf("Expected error response %v, but got %v", nil, errorResp)
		}
	})
}
