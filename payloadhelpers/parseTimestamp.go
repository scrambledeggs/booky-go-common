package payloadhelpers

import (
	"strings"
	"time"
)

func ParseTimestamp(datetime string) (time.Time, error) {
	formats := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02+15:04:05.999999999",
	}

	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		return time.Time{}, err
	}

	isUTC := strings.HasSuffix(datetime, "Z")
	for _, format := range formats {
		var parsedTime time.Time
		if isUTC {
			parsedTime, err = time.Parse(format, datetime)
		} else {
			parsedTime, err = time.ParseInLocation(format, datetime, loc)
		}

		if err == nil {
			return parsedTime.In(loc), nil
		}
	}

	return time.Time{}, err
}
