package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/communication/httptools"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestHealth(t *testing.T) {
	logger := logging.NewNopLogger()

	cfg := NewConfig(logger)
	cfg.Env = config.TestEnv

	db, err := postgres.Connect(
		logger,
		cfg.DBDsn,
		25,
		"15m",
		30,
		30*time.Second,
		5*time.Minute,
	)
	if err != nil {
		panic(err)
	}

	app := NewApp(logger, cfg, db)

	tReq := test.CreateRequestTester(
		app.Routes(),
		http.MethodGet,
		"/health",
	)
	rs := tReq.Do(t)

	var rsData Health
	err = httptools.ReadJSON(rs.Body, &rsData)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, true, rsData.IsDatabaseActive)
}
