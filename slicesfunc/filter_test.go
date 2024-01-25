package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	data := []string{"the", "quick", "brown", "potato", "sonof", "the", "potato", "the", "fudge"}

	newSlice := Filter(data, func(item string) bool {
		return item == "the"
	})

	assert.Equal(t, len(newSlice), 3)
}
