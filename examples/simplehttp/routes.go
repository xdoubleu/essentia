package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/pkg/middleware"
)

func (app Application) Routes() http.Handler {
	mux := http.NewServeMux()

	app.healthRoutes(mux)

	middleware, err := middleware.Default(
		app.logger,
		app.config.AllowedOrigins,
	)
	if err != nil {
		panic(err)
	}

	standard := alice.New(middleware...)
	return standard.Then(mux)
}
