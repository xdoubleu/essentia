package test_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/v2/pkg/communication/wstools"
	"github.com/xdoubleu/essentia/v2/pkg/logging"
	"github.com/xdoubleu/essentia/v2/pkg/test"
	"github.com/xdoubleu/essentia/v2/pkg/validate"
)

type TestResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type TestSubscribeMsg struct {
	TopicName string `json:"topicName"`
}

func (s TestSubscribeMsg) Validate() (bool, map[string]string) {
	v := validate.New()
	return v.Valid(), v.Errors()
}

func (s TestSubscribeMsg) Topic() string {
	return s.TopicName
}

func setup(t *testing.T) (http.Handler, *wstools.Topic) {
	t.Helper()

	logger := logging.NewNopLogger()
	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](
		logger,
		1,
		10,
	)
	topic, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		func(_ context.Context, _ *wstools.Topic) (any, error) {
			return TestResponse{Ok: true, Message: "initial"}, nil
		},
	)
	require.Nil(t, err)

	return ws.Handler(), topic
}

func TestWebSocketTester(t *testing.T) {
	handler, topic := setup(t)
	tWeb := test.CreateWebSocketTester(handler)
	tWeb.SetInitialMessage(TestSubscribeMsg{
		TopicName: "exists",
	})
	tWeb.SetParallelOperation(
		func(_ *testing.T, _ *httptest.Server) {
			topic.EnqueueEvent(TestResponse{
				Ok:      true,
				Message: "parallel",
			})
		},
	)

	var initialResponse, parallelResponse TestResponse
	err := tWeb.Do(t, &initialResponse, &parallelResponse)

	require.Nil(t, err)
	assert.Equal(t, TestResponse{Ok: true, Message: "initial"}, initialResponse)
	assert.Equal(t, TestResponse{Ok: true, Message: "parallel"}, parallelResponse)
}
