package httptools_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/internal/mocks"
	"github.com/xdoubleu/essentia/pkg/communication/httptools"
)

func TestServerFunc(t *testing.T) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         "localhost:8000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		time.Sleep(time.Second)
		err := srv.Shutdown(context.Background())
		require.Nil(t, err)
	}()

	mockedLogger := mocks.NewMockedLogger()
	err := httptools.Serve(mockedLogger.Logger(), srv, "test")

	require.Nil(t, err)
	assert.Contains(t, mockedLogger.CapturedLogs(), "starting")
	assert.Contains(t, mockedLogger.CapturedLogs(), "stopped")
}
