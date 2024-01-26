package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var filterData = []string{"the", "quick", "brown", "potato", "sonof", "the", "potato", "the", "fudge"}

func TestFilter(t *testing.T) {
	newSlice := Filter(filterData, func(item string) bool {
		return item == "the"
	})

	assert.Equal(t, len(newSlice), 3)
}

func ExampleFilter() {
	data := []string{"the", "quick", "brown", "potato", "sonof", "the", "potato", "the", "fudge"}
	newSlice := Filter(data, func(item string) bool {
		return item == "the"
	})

	fmt.Println(newSlice)

	// Output: [the the the]
}
