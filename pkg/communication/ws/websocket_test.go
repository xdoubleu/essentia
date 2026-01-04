package ws_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/test"
	"github.com/xdoubleu/essentia/pkg/validate"
)

type TestResponse struct {
	Ok bool `json:"ok"`
}

type TestSubscribeMsg struct {
	TopicName string `json:"topicName"`
}

func (s TestSubscribeMsg) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "topicName", s.TopicName, validate.IsNotEmpty)

	return v.Valid(), v.Errors()
}

func (s TestSubscribeMsg) Topic() string {
	return s.TopicName
}

func setup(t *testing.T) http.Handler {
	t.Helper()

	logger := logging.NewNopLogger()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](
		logger,
		1,
		10,
	)

	_, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		func(_ context.Context, _ *wstools.Topic) (any, error) {
			return TestResponse{Ok: true}, nil
		},
	)
	require.Nil(t, err)

	return ws.Handler()
}

func TestWebSocketExistingTopic(t *testing.T) {
	wsHandler := setup(t)

	tWeb := test.CreateWebSocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: "exists"})

	var initialResponse TestResponse
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.True(t, initialResponse.Ok)
}

func TestWebSocketUnknownTopic(t *testing.T) {
	wsHandler := setup(t)

	tWeb := test.CreateWebSocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: "unknown"})

	var initialResponse errortools.ErrorDto
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), initialResponse.Error)
	assert.Equal(t, "topic 'unknown' doesn't exist", initialResponse.Message)
}

func TestWebSocketExistingHandler(t *testing.T) {
	logger := logging.NewNopLogger()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](logger, 1, 10)
	topic, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		nil,
	)
	require.NotNil(t, topic)
	require.Nil(t, err)

	topic, err = ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		nil,
	)
	assert.Nil(t, topic)
	assert.EqualError(t, err, "topic 'exists' has already been added")
}

func TestWebsocketBasic(t *testing.T) {
	logger := logging.NewNopLogger()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](logger, 1, 10)
	topic, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		func(_ context.Context, _ *wstools.Topic) (any, error) {
			return TestResponse{Ok: true}, nil
		},
	)
	require.NotNil(t, topic)
	require.Nil(t, err)

	tWeb := test.CreateWebSocketTester(ws.Handler())
	tWeb.SetInitialMessage(TestSubscribeMsg{
		TopicName: "exists",
	})
	var rsData TestResponse
	err = tWeb.Do(t, &rsData, nil)
	require.Nil(t, err)

	assert.Equal(t, true, rsData.Ok)

	topic, err = ws.UpdateTopicName(topic, "new")
	require.NotNil(t, topic)
	require.Nil(t, err)

	tWeb = test.CreateWebSocketTester(ws.Handler())
	tWeb.SetInitialMessage(TestSubscribeMsg{
		TopicName: "new",
	})
	err = tWeb.Do(t, &rsData, nil)
	require.Nil(t, err)

	assert.Equal(t, true, rsData.Ok)

	err = ws.RemoveTopic(topic)
	assert.Nil(t, err)
}

func TestWebSocketUpdateExistingTopic(t *testing.T) {
	logger := logging.NewNopLogger()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](logger, 1, 10)
	topic, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		func(_ context.Context, _ *wstools.Topic) (any, error) {
			return TestResponse{Ok: true}, nil
		},
	)
	require.NotNil(t, topic)
	require.Nil(t, err)

	topic, err = ws.UpdateTopicName(topic, "exists")
	assert.Nil(t, topic)
	assert.ErrorContains(t, err, "topic 'exists' already exists")
}

func TestWebSocketUpdateNonExistingTopic(t *testing.T) {
	logger := logging.NewNopLogger()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](logger, 1, 10)
	topic, err := ws.UpdateTopicName(&wstools.Topic{
		Name: "unknown",
	}, "exists")
	assert.Nil(t, topic)
	assert.ErrorContains(t, err, "topic 'unknown' doesn't exist")
}

func TestWebSocketRemoveNonExistingTopic(t *testing.T) {
	logger := logging.NewNopLogger()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](logger, 1, 10)
	err := ws.RemoveTopic(&wstools.Topic{
		Name: "unknown",
	})
	assert.ErrorContains(t, err, "topic 'unknown' doesn't exist")
}
