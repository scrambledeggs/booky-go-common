package payloadhelpers

import (
	"strings"
	"time"

	"github.com/scrambledeggs/booky-go-common/converters"
)

func ParseTimestamp(datetime string) (date time.Time, err error) {
	formats := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02+15:04:05.999999999",
	}

	isUTC := strings.HasSuffix(datetime, "Z")
	for _, format := range formats {
		if isUTC {
			date, err = time.Parse(format, datetime)
		} else {
			loc, _ := time.LoadLocation("Asia/Manila")
			date, err = time.ParseInLocation(format, datetime, loc)
		}

		if err == nil {
			return converters.StringToTimeInTimezone(date.Format(time.RFC3339), "Asia/Manila")
		}
	}

	return
}
