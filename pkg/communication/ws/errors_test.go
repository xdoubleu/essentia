package ws_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	wstools "github.com/XDoubleU/essentia/pkg/communication/ws"
	"github.com/XDoubleU/essentia/pkg/config"
	errortools "github.com/XDoubleU/essentia/pkg/errors"
	"github.com/XDoubleU/essentia/pkg/logging"
	sentrytools "github.com/XDoubleU/essentia/pkg/sentry"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testErrorStatusCode(t *testing.T, handler http.HandlerFunc) int {
	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	res := httptest.NewRecorder()

	sentryMiddleware, err := sentrytools.Middleware(
		config.TestEnv,
		sentrytools.MockedSentryClientOptions(),
	)
	require.Nil(t, err)

	sentryMiddleware(handler).ServeHTTP(res, req)

	return res.Result().StatusCode
}

func setupWS(t *testing.T, allowedOrigin string) http.Handler {
	t.Helper()

	logger := logging.NewNopLogger()

	wsHandler := wstools.CreateWebSocketHandler[TestSubscribeMsg](
		logger,
		1,
		10,
	)
	_, err := wsHandler.AddTopic("topic", []string{allowedOrigin}, nil)
	require.Nil(t, err)

	sentryMiddleware, err := sentrytools.Middleware(
		config.TestEnv,
		sentrytools.MockedSentryClientOptions(),
	)
	require.Nil(t, err)

	return sentryMiddleware(wsHandler.Handler())
}

func TestUpgradeErrorResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		wstools.UpgradeErrorResponse(w, r, errors.New("test"))
	}

	statusCode := testErrorStatusCode(t, handler)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestErrorResponse(t *testing.T) {
	handler := setupWS(t, "http://localhost")

	tWeb := test.CreateWebSocketTester(handler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: "unknown"})

	var errorDto errortools.ErrorDto
	err := tWeb.Do(t, &errorDto, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, errorDto.Status)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), errorDto.Error)
	assert.Equal(t, "topic 'unknown' doesn't exist", errorDto.Message)
}

func TestFailedValidationResponse(t *testing.T) {
	handler := setupWS(t, "http://localhost")

	tWeb := test.CreateWebSocketTester(handler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: ""})

	var errorDto errortools.ErrorDto
	err := tWeb.Do(t, &errorDto, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, errorDto.Status)
	assert.Equal(t, http.StatusText(http.StatusUnprocessableEntity), errorDto.Error)
	assert.Equal(t, map[string]any{
		"topicName": "must be provided",
	}, errorDto.Message)
}
