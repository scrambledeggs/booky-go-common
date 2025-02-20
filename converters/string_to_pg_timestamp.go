package converters

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

func StringToPgTimestamp(datetime string) (date pgtype.Timestamptz, err error) {
	formats := []string{
		time.RFC3339, // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02+15:04:05.999999999",
	}

	var parsedTime time.Time
	for _, format := range formats {
		parsedTime, err = time.Parse(format, datetime)

		if err == nil {
			err = date.Scan(parsedTime.UTC())
			return
		}
	}

	date.Status = pgtype.Null
	return
}
