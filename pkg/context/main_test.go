package context_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	contexttools "github.com/xdoubleu/essentia/pkg/context"
	"github.com/xdoubleu/essentia/pkg/logging"
)

const testContextKey = contexttools.Key("test")

func TestGetContextValue(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testContextKey,
		true,
	)

	value := contexttools.GetValue[bool](ctx, testContextKey)

	assert.Equal(t, true, *value)
}

func TestGetContextValueNotPresent(t *testing.T) {
	value := contexttools.GetValue[bool](context.Background(), testContextKey)

	assert.Nil(t, value)
}

func TestGetContextValueIncorrectType(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testContextKey,
		10,
	)

	value := contexttools.GetValue[bool](ctx, testContextKey)

	assert.Nil(t, value)
}

func TestSetGetLogger(t *testing.T) {
	ctx := context.Background()

	logger := slog.Default()
	ctx = contexttools.WithLogger(ctx, logger)

	value := contexttools.Logger(ctx)

	assert.Equal(t, logger, value)
}

func TestGetNullLogger(t *testing.T) {
	ctx := context.Background()

	value := contexttools.Logger(ctx)

	assert.Equal(t, logging.NewNopLogger(), value)
}
