package converters

import (
	"time"
)

func StringToTimeInTimezone(datetime, timezone string) (time.Time, error) {
	utcTime, err := StringToTime(datetime)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return utcTime.In(location), nil
}
