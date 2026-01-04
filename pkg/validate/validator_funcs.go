package validate

import (
	"fmt"
	"slices"
	"time"

	"github.com/xdoubleu/essentia/v2/internal/helpers"
)

// ValidatorFunc is the expected format used for validating data using [Check].
type ValidatorFunc[T any] func(value T) (bool, string)

// IsNotEmpty checks that the provided value is not empty.
func IsNotEmpty(value string) (bool, string) {
	return value != "", "must be provided"
}

// IsGreaterThan checks if the provided value2 > value1.
func IsGreaterThan[T helpers.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 > value1, fmt.Sprintf("must be greater than %d", value1)
	}
}

// IsGreaterThanOrEqual checks if the provided value2 >= value1.
func IsGreaterThanOrEqual[T helpers.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 >= value1,
			fmt.Sprintf(
				"must be greater than or equal to %d",
				value1,
			)
	}
}

// IsLesserThan checks if the provided value2 < value1.
func IsLesserThan[T helpers.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 < value1, fmt.Sprintf("must be lesser than %d", value1)
	}
}

// IsLesserThanOrEqual checks if the provided value2 <= value1.
func IsLesserThanOrEqual[T helpers.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 <= value1,
			fmt.Sprintf(
				"must be lesser than or equal to %d",
				value1,
			)
	}
}

// IsValidTimeZone checks if the provided value is a valid IANA timezone.
func IsValidTimeZone(value string) (bool, string) {
	_, err := time.LoadLocation(value)
	return err == nil, "must be a valid IANA value"
}

// IsInSlice checks if the provided value is part of the provided slice.
func IsInSlice[T comparable](slice []T) ValidatorFunc[T] {
	return func(value T) (bool, string) {
		return slices.Contains(slice, value), "must be a valid value"
	}
}
