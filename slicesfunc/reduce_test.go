package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var reduceData = []string{"the", "quick", "brown"}

func TestReduce(t *testing.T) {
	reducedSlice := Reduce(reduceData, func(list map[string]string, item string) map[string]string {
		list[item] = fmt.Sprintf("%s-naknang", item)

		return list
	}, map[string]string{})

	assert.Equal(t, "the-naknang", reducedSlice["the"])
	assert.Equal(t, "quick-naknang", reducedSlice["quick"])
	assert.Equal(t, "brown-naknang", reducedSlice["brown"])
}

func ExampleReduce() {
	data := []string{"the", "quick", "brown"}
	reducedSlice := Reduce(data, func(list map[string]string, item string) map[string]string {
		list[item] = fmt.Sprintf("%s-naknang", item)

		return list
	}, map[string]string{})

	fmt.Println(reducedSlice)

	// Output: map[brown:brown-naknang quick:quick-naknang the:the-naknang]
}
