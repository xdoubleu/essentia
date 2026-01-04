package main

import (
	"context"
	"net/http"
	"time"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
)

func (app *Application) healthRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /health",
		app.getHealthHandler,
	)
}

type Health struct {
	IsDatabaseActive bool
}

func (app *Application) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	data := Health{
		IsDatabaseActive: app.db.Ping(ctx) == nil,
	}

	err := httptools.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
