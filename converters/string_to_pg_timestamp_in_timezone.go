package converters

import (
	"time"
)

func StringToTimeInTimezone(datetime string, timezone string) (date time.Time, err error) {
	utcTime, err := StringToTime(datetime)
	if err != nil {
		return date, err
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return date, err
	}

	date = utcTime.In(location)

	return
}
