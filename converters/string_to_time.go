package converters

import (
	"time"
)

func StringToTime(sDatetime string) (time.Time, error) {
	formats := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02+15:04:05.999999999",
	}

	var lastErr error
	for _, format := range formats {
		if datetime, err := time.Parse(format, sDatetime); err == nil {
			return datetime, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, lastErr
}
