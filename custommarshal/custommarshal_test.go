package custommarshal

import (
	"testing"

	"github.com/scrambledeggs/booky-brand-csmr/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestCustomMarshalMap(t *testing.T) {
	description := ""
	brand := entities.CombinedData{
		BrandData: entities.BrandData{
			BrandKey: entities.BrandKey{BrandID: 1},
			Brand: entities.Brand{
				Name:        "Renamed Brand 1",
				Status:      "inactive",
				Description: &description,
			},
		},
	}

	res, err := CustomMarshalMap(brand)
	assert.Equal(t, true, *res["description"].NULL)
	assert.Nil(t, err)
}
