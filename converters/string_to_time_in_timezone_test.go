package converters

import (
	"testing"
	"time"
)

func TestStringToTimeInTimezone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		timezone string
		expected time.Time
	}{
		{
			name:     "UTC to Asia/Manila",
			input:    "2025-02-20T15:04:05Z",
			timezone: "Asia/Manila",
			expected: time.Date(2025, 2, 20, 23, 4, 5, 0, time.FixedZone("Asia/Manila", 8*60*60)),
		},
		{
			name:     "New York to Asia/Manila",
			input:    "2025-02-20T15:04:05-05:00",
			timezone: "Asia/Manila",
			expected: time.Date(2025, 2, 21, 4, 4, 5, 0, time.FixedZone("Asia/Manila", 8*60*60)),
		},
		{
			name:     "Asia/Manila to Asia/Manila",
			input:    "2025-02-20T15:04:05+08:00",
			timezone: "Asia/Manila",
			expected: time.Date(2025, 2, 20, 15, 4, 5, 0, time.FixedZone("Asia/Manila", 8*60*60)),
		},
		{
			name:     "UTC to UTC",
			input:    "2025-02-20 15:04:05",
			timezone: "UTC",
			expected: time.Date(2025, 2, 20, 15, 4, 5, 0, time.UTC),
		},
		{
			name:     "Asia/Manila to UTC",
			input:    "2025-02-20T15:04:05+08:00",
			timezone: "UTC",
			expected: time.Date(2025, 2, 20, 7, 4, 5, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, err := StringToTimeInTimezone(tt.input, tt.timezone)
			if err != nil {
				t.Fatalf("[%s] unexpected error: %v", tt.name, err)
			}

			// Comparing in UTC for consistency
			if !date.UTC().Equal(tt.expected.UTC()) {
				t.Errorf("[%s] failed: expected %v, got %v", tt.name, tt.expected.UTC(), date.UTC())
			}
		})
	}
}

func TestStringToTimeInTimezone_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		timezone string
	}{
		{
			name:     "Invalid datetime",
			input:    "this is an invalid datetime",
			timezone: "Asia/Manila",
		},
		{
			name:     "Invalid date (Feb 30)",
			input:    "2025-02-30",
			timezone: "Asia/Manila",
		},
		{
			name:     "Invalid month (13)",
			input:    "2025-13-01",
			timezone: "Asia/Manila",
		},
		{
			name:     "Invalid time (99 minutes)",
			input:    "2025-02-20T15:99:99Z",
			timezone: "Asia/Manila",
		},
		{
			name:     "Invalid timezone",
			input:    "2025-02-20 15:04:05",
			timezone: "Invalid/Timezone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StringToTimeInTimezone(tt.input, tt.timezone)
			if err == nil {
				t.Errorf("Test %s failed: expected error for input %s with timezone %s", tt.name, tt.input, tt.timezone)
			}
		})
	}
}
