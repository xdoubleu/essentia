package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/justinas/alice"
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
