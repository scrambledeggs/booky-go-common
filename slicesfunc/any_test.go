package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	data := []string{"naknang", "patatas", "sonof", "potato"}

	existing := Any(data, func(s string) bool {
		return "naknang" == s
	})
	assert.Equal(t, true, existing)

	nonExisting := Any(data, func(s string) bool {
		return "bwakanang" == s
	})
	assert.Equal(t, false, nonExisting)
}
