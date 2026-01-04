package sentrytools

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// GoRoutineWrapper wraps a go routine with
// Sentry logic for error and performance tracking.
func GoRoutineWrapper(
	ctx context.Context,
	logger *slog.Logger,
	name string,
	f func(ctx context.Context, logger *slog.Logger) error,
) {
	name = fmt.Sprintf("GO ROUTINE %s", name)

	hub := sentry.CurrentHub().Clone()
	ctx = sentry.SetHubOnContext(ctx, hub)

	options := []sentry.SpanOption{
		sentry.WithOpName("go.routine"),
	}

	transaction := sentry.StartTransaction(ctx, name, options...)
	transaction.Status = sentry.HTTPtoSpanStatus(http.StatusOK)
	defer transaction.Finish()

	err := f(transaction.Context(), logger)

	if err != nil {
		transaction.Status = sentry.HTTPtoSpanStatus(http.StatusInternalServerError)

		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}
