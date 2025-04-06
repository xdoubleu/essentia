// Package grapher provides a tool for easily creating graphs.
package grapher

import (
	"fmt"
	"slices"
	"strconv"
	"time"
)

// GraphType is used to define the type of graph you need.
type GraphType int

// InterpolationType is used to define the type of interpolation you need.
type InterpolationType int

// Numeric is a union of all numeric types.
type Numeric interface {
	int | int64 | float64
}

const (
	// Normal GraphType provides a grapher where you can set points a certain dates.
	Normal GraphType = iota
	// Cumulative GraphType will accumulate all previous values on your graph.
	Cumulative GraphType = iota
	// CumulativeSameDate GraphType will accumulate values
	// with the same date on your graph.
	CumulativeSameDate GraphType = iota
)

const (
	// None InterpolationType provides a grapher that won't interpolate.
	None InterpolationType = iota
	// Zero InterpolationType provides a grapher that will interpolate using value zero.
	Zero InterpolationType = iota
	// PreviousValue InterpolationType provides a grapher
	// that will interpolate using the previous value.
	PreviousValue InterpolationType = iota
)

// Grapher is used to easily create graphs.
type Grapher[T Numeric] struct {
	graphType         GraphType
	interpolationType InterpolationType
	dateFormat        string
	dateGranularity   time.Duration
	dateStrings       []string
	values            map[string][]T
}

// New returns a new Grapher.
func New[T Numeric](
	graphType GraphType,
	interpolationType InterpolationType,
	dateFormat string,
	dateGranularity time.Duration,
) *Grapher[T] {
	return &Grapher[T]{
		graphType:         graphType,
		interpolationType: interpolationType,
		dateFormat:        dateFormat,
		dateGranularity:   dateGranularity,
		dateStrings:       []string{},
		values:            make(map[string][]T),
	}
}

// AddPoint adds a new point to the graph.
func (grapher *Grapher[T]) AddPoint(date time.Time, value T, label string) {
	dateStr := date.Format(grapher.dateFormat)

	dateIndex := grapher.getDateIndex(dateStr, label)

	grapher.updateDays(dateIndex, value, label)
}

func (grapher *Grapher[T]) getDateIndex(dateStr string, label string) int {
	dateIndex := slices.Index(grapher.dateStrings, dateStr)

	if dateIndex == -1 {
		grapher.addDates(dateStr)
		dateIndex = slices.Index(grapher.dateStrings, dateStr)
	}

	if _, ok := grapher.values[label]; !ok {
		grapher.values[label] = append(
			grapher.values[label],
			*new(T),
		)
	}

	for len(grapher.values[label]) < len(grapher.dateStrings) {
		grapher.values[label] = append(
			grapher.values[label],
			*new(T),
		)
	}

	return dateIndex
}

func (grapher *Grapher[T]) addDates(dateStr string) {
	if len(grapher.dateStrings) == 0 {
		grapher.dateStrings = append(grapher.dateStrings, dateStr)
		return
	}

	dateDay, _ := time.Parse(grapher.dateFormat, dateStr)
	smallestDate, _ := time.Parse(grapher.dateFormat, grapher.dateStrings[0])
	largestDate, _ := time.Parse(
		grapher.dateFormat,
		grapher.dateStrings[len(grapher.dateStrings)-1],
	)

	if grapher.interpolationType == None {
		if dateDay.Before(smallestDate) {
			grapher.addDateBefore(dateDay, *new(T))
		} else if dateDay.After(largestDate) {
			grapher.addDateAfter(dateDay)
		}
		return
	}

	i := smallestDate
	for i.After(dateDay) {
		i = i.Add(-1 * grapher.dateGranularity)
		grapher.addDateBefore(i, *new(T))
	}

	i = largestDate
	for i.Before(dateDay) {
		i = i.Add(grapher.dateGranularity)
		grapher.addDateAfter(i)
	}
}

func (grapher *Grapher[T]) addDateBefore(date time.Time, value T) {
	grapher.dateStrings = append(
		[]string{date.Format(grapher.dateFormat)},
		grapher.dateStrings...)

	for label := range grapher.values {
		grapher.values[label] = append(
			[]T{value},
			grapher.values[label]...)
	}
}

func (grapher *Grapher[T]) addDateAfter(date time.Time) {
	grapher.dateStrings = append(
		grapher.dateStrings,
		date.Format(grapher.dateFormat),
	)

	indexOfI := slices.Index(
		grapher.dateStrings,
		date.Format(grapher.dateFormat),
	)

	for label := range grapher.values {
		value := grapher.values[label][indexOfI-1]
		if grapher.interpolationType == Zero || grapher.interpolationType == None {
			value = *new(T)
		}

		grapher.values[label] = append(
			grapher.values[label],
			value,
		)
	}
}

func (grapher *Grapher[T]) updateDays(dateIndex int, value T, label string) {
	switch grapher.graphType {
	case Normal:
		grapher.values[label][dateIndex] = value
	case CumulativeSameDate:
		grapher.values[label][dateIndex] += value
	case Cumulative:
		for i := dateIndex; i < len(grapher.dateStrings); i++ {
			grapher.values[label][i] += value
		}
	}
}

// ToStringSlices returns the graph as a string slice of dates and values.
func (grapher Grapher[T]) ToStringSlices() ([]string, map[string][]string) {
	strValues := make(map[string][]string)

	for label, values := range grapher.values {
		for _, value := range values {
			strValue := ""
			switch v := any(value).(type) {
			case int:
				strValue = strconv.Itoa(v)
			case int64:
				strValue = strconv.Itoa(int(v))
			case float64:
				strValue = fmt.Sprintf("%.2f", v)
			}

			strValues[label] = append(strValues[label], strValue)
		}
	}

	return grapher.dateStrings, strValues
}

// ToSlices returns the graph as a string slice of dates and a typed slice of values.
func (grapher Grapher[T]) ToSlices() ([]string, map[string][]T) {
	return grapher.dateStrings, grapher.values
}
