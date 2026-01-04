package timetools_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/pkg/timetools"
)

func TestStartOfDay(t *testing.T) {
	now := time.Now()
	startOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	assert.Equal(t, startOfDay, timetools.StartOfDay(now))
}

func TestEndOfDay(t *testing.T) {
	now := time.Now()
	endOfDay := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		23, 59, 59, 999999999,
		now.Location(),
	)

	assert.Equal(t, endOfDay, timetools.EndOfDay(now))
}

func TestLocationIndependentTime(t *testing.T) {
	now := time.Now()

	utcTimeZone, _ := time.LoadLocation("UTC")
	result := timetools.LocationIndependentTime(now, now.Location().String())

	expected := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		utcTimeZone,
	)

	assert.Equal(t, expected, result)
}
