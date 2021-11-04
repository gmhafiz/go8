package time

import "time"

func Parse(date string, format ...string) time.Time {
	if format[0] != time.RFC3339 {
		panic("")
	} else if format[0] == time.RFC3339 {
		return parse3339(date)
	}
	panic("time format not supported")
}

func parse3339(rfc3339 string) time.Time {
	timeWant, _ := time.Parse(time.RFC3339, rfc3339)
	return timeWant
}
