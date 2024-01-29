package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var anyData = []string{"naknang", "patatas", "sonof", "potato"}

func TestAny(t *testing.T) {
	existing := Any(anyData, func(s string) bool {
		return s == "naknang"
	})

	assert.Equal(t, true, existing)
}

func TestAnyNonExisting(t *testing.T) {
	nonExisting := Any(anyData, func(s string) bool {
		return s == "bwakanang"
	})
	assert.Equal(t, false, nonExisting)
}

func ExampleAny() {
	list := []string{"naknang", "patatas", "sonof", "potato"}

	any := Any(list, func(s string) bool {
		return s == "naknang"
	})

	fmt.Println(any)

	// Output: true
}
