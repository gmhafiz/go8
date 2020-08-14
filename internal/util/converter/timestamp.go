package converter

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func TimestampToString(timestamp *timestamp.Timestamp) string {
	unixTimeUTC := time.Unix(timestamp.Seconds, 0)
	unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)

	return unitTimeInRFC3339
}

func StringToTime(timeInput string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, timeInput)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}
