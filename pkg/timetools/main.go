// Package timetools provides helpers for dealing with time.
package timetools

import "time"

// StartOfDay returns the start of day of the specified date.
func StartOfDay(dateTime time.Time) time.Time {
	output := time.Date(
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		0, 0, 0, 0,
		dateTime.Location(),
	)

	return output
}

// EndOfDay returns the end of day of the specified date.
func EndOfDay(dateTime time.Time) time.Time {
	output := time.Date(
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		23, 59, 59, 999999999,
		dateTime.Location(),
	)

	return output
}

// LocationIndependentTime returns the provided time in the
// provided time zone but forces the time zone to UTC.
func LocationIndependentTime(t time.Time, locationTimeZone string) time.Time {
	timeZone, _ := time.LoadLocation(locationTimeZone)
	utcTimeZone, _ := time.LoadLocation("UTC")
	t = t.In(timeZone)
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		utcTimeZone,
	)
}
