package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var findData = []string{"the", "quick", "brown", "potato"}

func TestFind(t *testing.T) {
	brown, ok := Find(findData, func(item string) bool {
		return item == "brown"
	})
	assert.Equal(t, "brown", *brown)
	assert.Equal(t, true, ok)
}

func TestFindNonExisting(t *testing.T) {
	_, ok := Find(findData, func(item string) bool {
		return item == "naknang"
	})
	assert.Equal(t, false, ok)
}
