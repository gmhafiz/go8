package converter

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jinzhu/now"
)

func TimestampToString(timestamp *timestamp.Timestamp) string {
	unixTimeUTC := time.Unix(timestamp.Seconds, 0)
	unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339)

	return unitTimeInRFC3339
}

func StringToTime(converter now.Config, timeInput string) (time.Time, error) {
	parsed, err := converter.Parse(timeInput)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}