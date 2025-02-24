package payloadhelpers

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "UTC timestamp with Z",
			input:    "2024-02-21T12:00:00Z",
			expected: "2024-02-21T20:00:00+08:00",
		},
		{
			name:     "Timestamp with explicit +08:00 timezone",
			input:    "2024-02-21T12:00:00+08:00",
			expected: "2024-02-21T12:00:00+08:00",
		},
		{
			name:     "Timestamp with explicit -05:00 timezone",
			input:    "2024-02-21T12:00:00-05:00",
			expected: "2024-02-22T01:00:00+08:00",
		},
		{
			name:     "Timestamp without timezone (should default to midnight in Asia/Manila)",
			input:    "2024-02-21T12:00:00",
			expected: "2024-02-21T12:00:00+08:00",
		},
		{
			name:     "Date-only format (should default to midnight in Asia/Manila)",
			input:    "2024-02-21",
			expected: "2024-02-21T00:00:00+08:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseTimestamp(tt.input)
			if err != nil {
				t.Errorf("[%s] Failed to parse timestamp: %v", tt.name, err)
				return
			}

			// Reformat for easy testing
			result := parsed.Time.Format(time.RFC3339)

			if result != tt.expected {
				t.Errorf("[%s] failed: expected %v but got %v", tt.name, tt.expected, result)
			}
		})
	}
}
