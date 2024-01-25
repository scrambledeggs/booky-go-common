package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	data := []string{"naknang", "patatas", "sonof", "potato"}

	existing := Contains("naknang", data)
	assert.Equal(t, true, existing)

	nonExisting := Contains("bwakanang", data)
	assert.Equal(t, false, nonExisting)
}
