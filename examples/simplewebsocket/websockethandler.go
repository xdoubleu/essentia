package main

import (
	"context"
	"net/http"

	wstools "github.com/XDoubleU/essentia/pkg/communication/ws"
	"github.com/XDoubleU/essentia/pkg/validate"
)

type SubscribeMessageDto struct {
	TopicName string `json:"topicName"`
}

type ResponseMessageDto struct {
	Message string `json:"message"`
}

func (msg SubscribeMessageDto) Validate() (bool, map[string]string) {
	v := validate.New()
	return v.Valid(), v.Errors()
}

func (msg SubscribeMessageDto) Topic() string {
	return msg.TopicName
}

func (app *Application) websocketRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /",
		app.getWebSocketHandler(),
	)
}

func (app *Application) getWebSocketHandler() http.HandlerFunc {
	wsHandler := wstools.CreateWebSocketHandler[SubscribeMessageDto](
		app.logger,
		1,
		10, //nolint:mnd //no magic number
	)
	_, err := wsHandler.AddTopic(
		"topic",
		app.config.AllowedOrigins,
		func(_ context.Context, _ *wstools.Topic) (any, error) {
			return ResponseMessageDto{
				Message: "Hello, World!",
			}, nil
		})
	if err != nil {
		panic(err)
	}

	return wsHandler.Handler()
}
