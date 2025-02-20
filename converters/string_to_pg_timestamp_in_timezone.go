package converters

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

func StringToPgTimestampInTimezone(datetime string, timezone string) (date pgtype.Timestamptz, err error) {
	utcTime, err := StringToPgTimestamp(datetime)
	if err != nil {
		return date, err
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return date, err
	}

	localTime := utcTime.Time.In(location)

	err = date.Scan(localTime)

	return
}
