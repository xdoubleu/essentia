// Package errortools contains reusable error messages
// and other helpers for dealing with errors.
package errortools

//nolint:lll,revive //can't make these lines shorter,the errors are clear
const (
	MessageInternalServerError = "the server encountered a problem and could not process your request"
	MessageTooManyRequests     = "rate limit exceeded"
	MessageForbidden           = "user has no access to this resource"
)
