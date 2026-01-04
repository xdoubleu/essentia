package sentrytools_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/sentrytools"
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
