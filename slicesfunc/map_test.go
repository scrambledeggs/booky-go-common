package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	data := []string{"the", "quick", "brown"}

	newSlice := Map(data, func(item string) string {
		return fmt.Sprintf("%s-potato", item)
	})
	assert.Equal(t, newSlice[0], "the-potato")
	assert.Equal(t, newSlice[1], "quick-potato")
	assert.Equal(t, newSlice[2], "brown-potato")
}
