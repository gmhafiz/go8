package time

import (
	"log"
	"time"
)

func Parse(date string, format ...string) time.Time {
	if len(format) == 0 {
		if len(date) == 10 {
			return parseISO8601(date)
		} else {
			return parse3339(date)
		}
	} else if format[0] == time.RFC3339 {
		return parse3339(date)
	} else if len(format) == 1 {
		return parse3339(format[0])
	}
	panic("time format not supported")
}

func parseISO8601(iso8601 string) time.Time {
	timeWant, err := time.Parse("2006-01-02T15:04:05", iso8601)
	if err != nil {
		log.Panic(err)
	}
	return timeWant
}

func parse3339(rfc3339 string) time.Time {
	timeWant, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		log.Println(err)
	}
	return timeWant
}
