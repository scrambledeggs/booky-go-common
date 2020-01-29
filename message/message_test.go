package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	message, err := NewMessageEvent()
	assert.NotNil(t, message)
	assert.Nil(t, err)
}

func TestGuuid(t *testing.T) {
	message, err := NewMessageEvent()
	assert.NotNil(t, message)
	assert.NotNil(t, message.GetUUID())
	assert.Nil(t, err)
}
