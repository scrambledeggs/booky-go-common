package slicesfunc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var countData = []string{"the", "quick", "brown", "potato", "sonof", "the", "potato", "the", "fudge"}

func TestCount(t *testing.T) {
	theCount := Count(countData, "the")
	assert.Equal(t, 3, theCount)

	potatoCount := Count(countData, "potato")
	assert.Equal(t, 2, potatoCount)

	quickCount := Count(countData, "quick")
	assert.Equal(t, 1, quickCount)

	naknangCount := Count(countData, "naknang")
	assert.Equal(t, 0, naknangCount)
}

func ExampleCount() {
	data := []string{"the", "quick", "brown", "potato", "sonof", "the", "potato", "the", "fudge"}

	theCount := Count(data, "the")
	fmt.Println(theCount)

	// Output: 3
}
