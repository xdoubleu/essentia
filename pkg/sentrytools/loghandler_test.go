package sentrytools_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/v4/pkg/config"
	"github.com/xdoubleu/essentia/v4/pkg/logging"
	"github.com/xdoubleu/essentia/v4/pkg/sentrytools"
)

func TestLogHandlerDev(t *testing.T) {
	var buf bytes.Buffer

	logger := slog.New(
		sentrytools.NewLogHandler(
			config.DevEnv,
			//nolint:exhaustruct //other fields are optional
			logging.NewBufLogHandler(
				&buf,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	logger.Error("test", logging.ErrAttr(errors.New("testerror")))

	assert.Contains(t, buf.String(), "level=ERROR msg=test error=testerror")
}

func TestLogHandlerWith(t *testing.T) {
	var buf bytes.Buffer

	logger := slog.New(
		sentrytools.NewLogHandler(
			config.DevEnv,
			//nolint:exhaustruct //other fields are optional
			logging.NewBufLogHandler(
				&buf,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	logger = logger.With(slog.String("source", "test"))

	logger.Error("test", logging.ErrAttr(errors.New("testerror")))

	test := buf.String()
	assert.Contains(t, test, "level=ERROR msg=test source=test error=testerror")
}

func TestLogHandlerSentryReceivesActualError(t *testing.T) {
	hub := sentrytools.MockedSentryHub()
	ctx := sentry.SetHubOnContext(context.Background(), hub)

	var buf bytes.Buffer
	logger := slog.New(
		sentrytools.NewLogHandler(
			config.DevEnv,
			//nolint:exhaustruct //other fields are optional
			logging.NewBufLogHandler(
				&buf,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	sentErr := errors.New("actual sentinel error")
	logger.ErrorContext(ctx, "something went wrong", "error", sentErr, "key", "val")

	events := sentrytools.MockedHubEvents(hub)
	require.Len(t, events, 1)
	assert.Equal(t, "actual sentinel error", events[0].Exception[0].Value)
}
