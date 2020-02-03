package message

import (
	guuid "github.com/google/uuid"
)

// Message Interface
type Message interface {
	GetUUID() *guuid.UUID
}

// EventMessage struct
type EventMessage struct {
	guuid   *guuid.UUID
	payload interface{}
}

// NewEventMessage constructor
func NewEventMessage() (Message, error) {
	strGuuid, err := guuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &EventMessage{&strGuuid, nil}, nil
}

// GetUUID method
func (e *EventMessage) GetUUID() *guuid.UUID {
	return e.guuid
}
