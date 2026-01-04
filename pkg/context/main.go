// Package context provides functions which can be used to
// set/get values to/from [context.Context].
package context

import (
	"context"
	"log/slog"

	"github.com/xdoubleu/essentia/pkg/logging"
)

// Key is the type used for specifying context keys.
type Key string

// GetValue returns a value by key from the context.
func GetValue[T any](ctx context.Context, key Key) *T {
	val := ctx.Value(key)
	if val == nil {
		return nil
	}

	castedValue, ok := val.(T)
	if !ok {
		return nil
	}

	return &castedValue
}

// WithLogger sets the logger on the context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// Logger returns the logger stored in the context or a NopLogger.
func Logger(ctx context.Context) *slog.Logger {
	logger := GetValue[*slog.Logger](ctx, loggerContextKey)

	if logger == nil {
		return logging.NewNopLogger()
	}

	return *logger
}

// WithShownErrors enables showing errors
// of [httptools.ServerErrorResponse].
func WithShownErrors(ctx context.Context) context.Context {
	return context.WithValue(ctx, showErrorsContextKey, true)
}

// ShowErrors returns if errors should be shown
// in [httptools.ServerErrorResponse.].
func ShowErrors(ctx context.Context) bool {
	showErrors := GetValue[bool](ctx, showErrorsContextKey)

	if showErrors == nil {
		return false
	}

	return *showErrors
}
