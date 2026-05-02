package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/v3/pkg/middleware"
)

func (app Application) Routes() http.Handler {
	mux := http.NewServeMux()

	app.websocketRoutes(mux)

	middleware, err := middleware.DefaultWithSentry(
		app.logger,
		app.config.AllowedOrigins,
		app.config.Env,
	)
	if err != nil {
		panic(err)
	}

	standard := alice.New(middleware...)
	return standard.Then(mux)
}
