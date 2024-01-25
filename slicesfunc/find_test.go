package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	data := []string{"the", "quick", "brown", "potato"}

	brown, ok := Find(data, func(item string) bool {
		return item == "brown"
	})
	assert.Equal(t, "brown", *brown)
	assert.Equal(t, true, ok)

	_, ok = Find(data, func(item string) bool {
		return item == "naknang"
	})
	assert.Equal(t, false, ok)
}
