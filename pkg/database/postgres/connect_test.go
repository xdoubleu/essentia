package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/internal/mocks"
	"github.com/xdoubleu/essentia/v2/pkg/database/postgres"
)

func TestConnectRetries(t *testing.T) {
	mockedLogger := mocks.NewMockedLogger()

	_, err := postgres.Connect(
		mockedLogger.Logger(),
		"",
		1,
		"1s",
		1,
		200*time.Millisecond,
		time.Second,
	)

	assert.NotNil(t, err)
	assert.Contains(t, mockedLogger.CapturedLogs(), "retry_in")
}
