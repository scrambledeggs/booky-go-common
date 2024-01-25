package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	data := []string{"the", "quick", "brown"}

	reducedSlice := Reduce(data, func(list map[string]string, item string) map[string]string {
		list[item] = fmt.Sprintf("%s-naknang", item)

		return list
	}, map[string]string{})

	assert.Equal(t, "the-naknang", reducedSlice["the"])
	assert.Equal(t, "quick-naknang", reducedSlice["quick"])
	assert.Equal(t, "brown-naknang", reducedSlice["brown"])
}
