package slicesfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	data := []string{"the", "quick", "brown", "potato", "sonof", "the", "potato", "the", "fudge"}

	theCount := Count(data, "the")
	assert.Equal(t, 3, theCount)

	potatoCount := Count(data, "potato")
	assert.Equal(t, 2, potatoCount)

	quickCount := Count(data, "quick")
	assert.Equal(t, 1, quickCount)

	naknangCount := Count(data, "naknang")
	assert.Equal(t, 0, naknangCount)
}
