package main

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/v2/pkg/middleware"
)

func (app Application) Routes() http.Handler {
	mux := http.NewServeMux()

	app.websocketRoutes(mux)

	middleware, err := middleware.DefaultWithSentry(
		app.logger,
		app.config.AllowedOrigins,
		app.config.Env,
		//nolint:exhaustruct //not all fields are needed
		sentry.ClientOptions{},
	)
	if err != nil {
		panic(err)
	}

	standard := alice.New(middleware...)
	return standard.Then(mux)
}
