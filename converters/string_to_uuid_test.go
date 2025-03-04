package converters

import (
	"testing"
)

func TestStringToUUID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Valid UUID",
			input:       "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:        "Invalid UUID",
			input:       "not-a-valid-uuid",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, err := StringToUUID(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for input %q, but got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tt.input, err)
				}
				if uuid.String() != tt.input {
					t.Errorf("expected UUID %q, got %q", tt.input, uuid.String())
				}
			}
		})
	}
}
