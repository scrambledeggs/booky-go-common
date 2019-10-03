package marshalling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Offer struct {
	Title string `json:"title,omitempty"`
}

type Brand struct {
	Name        string  `json:"brand_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Kind        string  `json:"kind,omitempty"`
	Status      string  `json:"brand_status,omitempty"`
	Limit       *int    `json:"offer_limit,omitempty"`
	Offer       *Offer  `json:"offer"`
}

func TestCustomMarshalMap(t *testing.T) {
	description := ""
	zero := 0
	brand := Brand{
		Name:        "Renamed Brand 1",
		Status:      "inactive",
		Description: &description,
		Limit:       &zero,
		Offer:       &Offer{},
	}

	res, err := CustomMarshalMap(brand)
	assert.Equal(t, true, *res["description"].NULL)
	assert.Equal(t, "0", *res["offer_limit"].N)
	assert.Equal(t, true, *res["offer"].NULL)
	assert.Nil(t, err)
}
