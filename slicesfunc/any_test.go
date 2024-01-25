package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var anyData = []string{"naknang", "patatas", "sonof", "potato"}

func TestAny(t *testing.T) {
	existing := Any(anyData, func(s string) bool {
		return "naknang" == s
	})
	assert.Equal(t, true, existing)
}

func TestAnyNonExisting(t *testing.T) {
	nonExisting := Any(anyData, func(s string) bool {
		return "bwakanang" == s
	})
	assert.Equal(t, false, nonExisting)
}
