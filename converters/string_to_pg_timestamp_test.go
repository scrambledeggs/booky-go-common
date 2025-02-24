package converters

import (
	"testing"
	"time"
)

func TestStringToPgTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "Valid UTC RFC3339 format",
			input:    "2025-02-20T15:04:05Z",
			expected: time.Date(2025, 2, 20, 15, 4, 5, 0, time.UTC),
		},
		{
			name:     "Valid date only format in UTC",
			input:    "2025-02-20",
			expected: time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Valid date and time format in Asia/Manila (+08:00)",
			input:    "2025-02-20T15:04:05+08:00",
			expected: time.Date(2025, 2, 20, 7, 4, 5, 0, time.UTC),
		},
		{
			name:     "Valid date and time format in America/New_York (-05:00)",
			input:    "2025-02-20T15:04:05-05:00",
			expected: time.Date(2025, 2, 20, 20, 4, 5, 0, time.UTC),
		},
		{
			name:     "Valid date and time format in Europe/London (+00:00)",
			input:    "2025-02-20T15:04:05+00:00",
			expected: time.Date(2025, 2, 20, 15, 4, 5, 0, time.UTC),
		},
		{
			name:     "Valid date and time format in Australia/Sydney (+11:00)",
			input:    "2025-02-20T15:04:05+11:00",
			expected: time.Date(2025, 2, 20, 4, 4, 5, 0, time.UTC),
		},
		{
			name:     "Valid date and time with nanoseconds in Asia/Manila (+08:00)",
			input:    "2025-02-20T15:04:05.999999999+08:00",
			expected: time.Date(2025, 2, 20, 7, 4, 5, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, err := StringToPgTimestamp(tt.input)
			if err != nil {
				t.Fatalf("unexpected error in test %s: %v", tt.name, err)
			}

			if !date.Time.Equal(tt.expected) {
				t.Errorf("[%s] failed: expected %v, got %v", tt.name, tt.expected, date.Time)
			}
		})
	}
}

func TestStringToPgTimestamp_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "Invalid datetime", input: "invalid datetime"},
		{name: "Invalid date (Feb 30)", input: "2025-02-30"},
		{name: "Invalid month (13)", input: "2025-13-01"},
		{name: "Invalid time (99 minutes)", input: "2025-02-20T15:99:99Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StringToPgTimestamp(tt.input)
			if err == nil {
				t.Errorf("Test %s failed: expected error for input %s", tt.name, tt.input)
			}
		})
	}
}
