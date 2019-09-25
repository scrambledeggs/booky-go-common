package photo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatImageURL(t *testing.T) {
	idURL := FormatIDUrl(20)
	assert.Equal(t, "000/000/020", idURL)

	imageUrl := FormatImageURL(20, "listings", "sample.jpg")
	assert.Equal(t, "https://assets1.phonebooky.com/listings/photos/000/000/020/original/sample.jpg", imageUrl)

	imageUrl = FormatImageURL(20, "listings", "sample.jpg", "medium")
	assert.Equal(t, "https://assets1.phonebooky.com/listings/photos/000/000/020/medium/sample.jpg", imageUrl)

	imageUrl = FormatImageURL(20, "listings", "sample.jpg", "medium", "logos")
	assert.Equal(t, "https://assets1.phonebooky.com/listings/logos/000/000/020/medium/sample.jpg", imageUrl)

	imageUrl = FormatImageURL(20, "listings", "")
	assert.Equal(t, "", imageUrl)
}
