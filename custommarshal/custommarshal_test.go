package custommarshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Brand struct {
	Name        string  `json:"brand_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Kind        string  `json:"kind,omitempty"`
	Status      string  `json:"brand_status,omitempty"`
}

func TestMap(t *testing.T) {
	description := ""
	brand := Brand{
		Name:        "Renamed Brand 1",
		Status:      "inactive",
		Description: &description,
	}

	res, err := Map(brand)
	assert.Equal(t, true, *res["description"].NULL)
	assert.Nil(t, err)
}
