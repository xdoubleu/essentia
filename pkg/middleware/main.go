// Package middleware provides configurable middleware and predefined lists,
// such as [Minimal], [Default] and [DefaultWithSentry].
package middleware

import (
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/goddtriffin/helmet"
	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/pkg/sentrytools"
)

// Minimal provides a predefined chain of useful middleware.
// Being:
//   - [Logger]
//   - [Recover]
func Minimal(logger *slog.Logger) []alice.Constructor {
	return []alice.Constructor{
		Logger(logger),
		Recover(logger),
	}
}

// Default provides a predefined chain of useful middleware.
// Being:
//   - All middleware from [Minimal]
//   - [helmet.Helmet]
//   - [CORS]
//   - [RateLimit]
func Default(
	logger *slog.Logger,
	allowedOrigins []string,
) ([]alice.Constructor, error) {
	return defaultBase(logger, allowedOrigins, nil, nil)
}

// DefaultWithSentry provides a predefined chain of useful middleware.
// Being:
//   - All middleware from [Default]
//   - [sentrytools.Middleware]
func DefaultWithSentry(
	logger *slog.Logger,
	allowedOrigins []string,
	env string,
	sentryClientOptions sentry.ClientOptions,
) ([]alice.Constructor, error) {
	return defaultBase(logger, allowedOrigins, &env, &sentryClientOptions)
}

func defaultBase(
	logger *slog.Logger,
	allowedOrigins []string,
	env *string,
	sentryClientOptions *sentry.ClientOptions,
) ([]alice.Constructor, error) {
	useSentry := env != nil && sentryClientOptions != nil

	helmet := helmet.Default()

	handlers := Minimal(logger)
	handlers = append(handlers, helmet.Secure)
	handlers = append(handlers, CORS(allowedOrigins, true))
	//nolint:mnd//no magic number
	handlers = append(handlers, RateLimit(10, 30, time.Minute, 3*time.Minute))

	if useSentry {
		sentryMiddleware, err := sentrytools.Middleware(*env, *sentryClientOptions)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, sentryMiddleware)
	}

	return handlers, nil
}
