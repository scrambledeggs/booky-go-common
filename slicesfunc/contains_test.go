package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var containsData = []string{"naknang", "patatas", "sonof", "potato"}

func TestContains(t *testing.T) {
	existing := Contains("naknang", containsData)
	assert.Equal(t, true, existing)
}

func TestContainsNonExisting(t *testing.T) {
	nonExisting := Contains("bwakanang", containsData)
	assert.Equal(t, false, nonExisting)
}
