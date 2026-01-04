package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/logging"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"
)

type Application struct {
	logger *slog.Logger
	config Config
}

func NewApp(logger *slog.Logger, config Config) Application {
	return Application{
		logger: logger,
		config: config,
	}
}

func main() {
	cfg := NewConfig(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	logger := slog.New(
		sentrytools.NewLogHandler(cfg.Env, slog.NewTextHandler(os.Stdout, nil)),
	)

	app := NewApp(logger, cfg)

	//nolint:mnd //no magic number
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Error("failed to serve server", logging.ErrAttr(err))
	}
}
