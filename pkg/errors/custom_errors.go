package errors

import (
	"fmt"

	"github.com/xdoubleu/essentia/internal/shared"
)

// NotFoundError is used when a certain resource doesn't exist.
type NotFoundError struct {
	resourceName    string
	identifierValue string
	JSONField       string
}

// ConflictError is used when an existing resource would conflict with a new resource.
type ConflictError struct {
	resourceName    string
	identifierValue string
	JSONField       string
}

// BadRequestError is used to return a bad request response.
type BadRequestError struct {
	err error
}

// UnauthorizedError is used to return an unauthorized response.
type UnauthorizedError struct {
	err error
}

// NewNotFoundError creates a new [NotFoundError].
func NewNotFoundError(
	resourceName string,
	identifierValue any,
	jsonField string,
) NotFoundError {
	value, err := shared.AnyToString(identifierValue)
	if err != nil {
		panic(err)
	}

	return NotFoundError{
		resourceName:    resourceName,
		identifierValue: value,
		JSONField:       jsonField,
	}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf(
		"%s with %s '%s' doesn't exist",
		err.resourceName,
		err.JSONField,
		err.identifierValue,
	)
}

// NewConflictError creates a new [ConflictError].
func NewConflictError(
	resourceName string,
	identifierValue any,
	jsonField string,
) ConflictError {
	value, err := shared.AnyToString(identifierValue)
	if err != nil {
		panic(err)
	}

	return ConflictError{
		resourceName:    resourceName,
		identifierValue: value,
		JSONField:       jsonField,
	}
}

func (err ConflictError) Error() string {
	return fmt.Sprintf(
		"%s with %s '%s' already exists",
		err.resourceName,
		err.JSONField,
		err.identifierValue,
	)
}

// NewBadRequestError creates a new [BadRequestError].
func NewBadRequestError(err error) BadRequestError {
	return BadRequestError{
		err: err,
	}
}

func (err BadRequestError) Error() string {
	return err.err.Error()
}

// NewUnauthorizedError creates a new [UnauthorizedError].
func NewUnauthorizedError(err error) UnauthorizedError {
	return UnauthorizedError{
		err: err,
	}
}

func (err UnauthorizedError) Error() string {
	return err.err.Error()
}
