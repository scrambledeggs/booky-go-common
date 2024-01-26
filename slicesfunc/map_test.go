package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mapData = []string{"the", "quick", "brown"}

func TestMap(t *testing.T) {
	newSlice := Map(mapData, func(item string) string {
		return fmt.Sprintf("%s-potato", item)
	})
	assert.Equal(t, newSlice[0], "the-potato")
	assert.Equal(t, newSlice[1], "quick-potato")
	assert.Equal(t, newSlice[2], "brown-potato")
}

func ExampleMap() {
	data := []string{"the", "quick", "brown"}

	newSlice := Map(data, func(item string) string {
		return fmt.Sprintf("%s-potato", item)
	})

	fmt.Println(newSlice)

	// Output: [the-potato quick-potato brown-potato]
}
