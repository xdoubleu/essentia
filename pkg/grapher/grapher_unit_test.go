package grapher_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/grapher"
)

func TestGrapherCumulative(t *testing.T) {
	grapher := grapher.New[int](
		grapher.Cumulative,
		grapher.PreviousValue,
		"2006-01-02",
		24*time.Hour,
	)

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), 1, "data")
	}

	for i := 1; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, -1*i), 1, "data")
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 19, len(dateSlice))
	assert.Equal(t, 19, len(valueSlice["data"]))

	for i := 0; i < 19; i++ {
		assert.Equal(
			t,
			time.Now().UTC().AddDate(0, 0, i-9).Format("2006-01-02"),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i+1), valueSlice["data"][i])
	}
}

func TestGrapherNormal(t *testing.T) {
	grapher := grapher.New[int](
		grapher.Normal,
		grapher.PreviousValue,
		"2006-01-02",
		24*time.Hour,
	)

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), i, "data")
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 10, len(dateSlice))
	assert.Equal(t, 10, len(valueSlice["data"]))

	for i := 0; i < 10; i++ {
		assert.Equal(
			t,
			time.Now().UTC().AddDate(0, 0, i).Format("2006-01-02"),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i), valueSlice["data"][i])
	}
}

func TestGrapherNormalSeconds(t *testing.T) {
	grapher := grapher.New[int](grapher.Normal, grapher.None, time.RFC3339, time.Second)

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.Add(time.Duration(i)*time.Second), i, "data")
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 10, len(dateSlice))
	assert.Equal(t, 10, len(valueSlice["data"]))

	for i := 0; i < 10; i++ {
		assert.Equal(
			t,
			time.Now().UTC().Add(time.Duration(i)*time.Second).Format(time.RFC3339),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i), valueSlice["data"][i])
	}
}

func TestGrapherShiftedLabels(t *testing.T) {
	grapher := grapher.New[int](
		grapher.Normal,
		grapher.PreviousValue,
		"2006-01-02",
		24*time.Hour,
	)

	dateNow := time.Now().UTC()
	for i := 0; i < 2; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), i, "data1")
	}

	for i := 0; i < 2; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), i, "data2")
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 2, len(dateSlice))
	assert.Equal(t, 2, len(valueSlice["data1"]))
	assert.Equal(t, 2, len(valueSlice["data2"]))
}
