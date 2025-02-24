package converters

import (
	"time"
)

func StringToTime(sDatetime string) (datetime time.Time, err error) {
	formats := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02+15:04:05.999999999",
	}

	for _, format := range formats {
		datetime, err = time.Parse(format, sDatetime)

		if err == nil {
			return
		}
	}

	return
}
