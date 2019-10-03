package marshalling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Brand struct {
	Name        string  `json:"brand_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Kind        string  `json:"kind,omitempty"`
	Status      string  `json:"brand_status,omitempty"`
	Limit       *int    `json:"offer_limit,omitempty"`
}

func TestCustomMarshalMap(t *testing.T) {
	description := ""
	zero := 0
	brand := Brand{
		Name:        "Renamed Brand 1",
		Status:      "inactive",
		Description: &description,
		Limit:       &zero,
	}

	res, err := CustomMarshalMap(brand)
	assert.Equal(t, true, *res["description"].NULL)
	assert.Equal(t, "0", *res["offer_limit"].N)
	assert.Nil(t, err)
}
