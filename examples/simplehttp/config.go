package main

import (
	"log/slog"

	"github.com/XDoubleU/essentia/pkg/config"
)

type Config struct {
	Env            string
	Port           int
	DBDsn          string
	AllowedOrigins []string
}

func NewConfig(logger *slog.Logger) Config {
	c := config.New(logger)

	var cfg Config

	cfg.Env = c.EnvStr("ENV", config.ProdEnv)
	//nolint:mnd //no magic number
	cfg.Port = c.EnvInt("PORT", 8000)
	cfg.DBDsn = c.EnvStr("DB_DSN", "postgres://postgres@localhost/postgres")
	cfg.AllowedOrigins = c.EnvStrArray(
		"ALLOWED_ORIGINS",
		[]string{"http://localhost"},
	)

	return cfg
}
