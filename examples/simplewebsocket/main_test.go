package main

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestWebSocket(t *testing.T) {
	cfg := NewConfig(logging.NewNopLogger())
	cfg.Env = config.TestEnv

	logger := slog.New(
		sentrytools.NewLogHandler(cfg.Env, slog.NewTextHandler(os.Stdout, nil)),
	)
	app := NewApp(logger, cfg)

	tWeb := test.CreateWebSocketTester(app.Routes())
	tWeb.SetInitialMessage(SubscribeMessageDto{
		TopicName: "topic",
	})

	var initialResponse ResponseMessageDto
	err := tWeb.Do(t, &initialResponse, nil)
	require.Nil(t, err)

	assert.Equal(t, "Hello, World!", initialResponse.Message)
}
