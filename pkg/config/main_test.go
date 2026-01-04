package config_test

import (
	"bytes"
	"log/slog"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
)

const existingKey, nonExistingKey = "key", "non_key"

func TestEnvStr(t *testing.T) {
	expected, def := "string", ""

	t.Setenv(existingKey, expected)

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	exists := c.EnvStr(existingKey, def)
	notExists := c.EnvStr(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)

	assert.Contains(t, buf.String(), "loaded env var 'key'='string' with type 'string'")
	assert.Contains(t, buf.String(), "loaded env var 'non_key'='' with type 'string'")
}

func TestEnvStrArray(t *testing.T) {
	rawExpected := "string1,string2"
	expected, def := []string{"string1", "string2"}, []string{""}

	t.Setenv(existingKey, rawExpected)

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	exists := c.EnvStrArray(existingKey, def)
	notExists := c.EnvStrArray(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)

	assert.Contains(
		t,
		buf.String(),
		"loaded env var 'key'='string1,string2' with type 'string array'",
	)
	assert.Contains(
		t,
		buf.String(),
		"loaded env var 'non_key'='' with type 'string array'",
	)
}

func TestEnvInt(t *testing.T) {
	expected, def := 14, 0

	t.Setenv(existingKey, strconv.Itoa(expected))

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	exists := c.EnvInt(existingKey, def)
	notExists := c.EnvInt(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)

	assert.Contains(t, buf.String(), "loaded env var 'key'='14' with type 'int'")
	assert.Contains(t, buf.String(), "loaded env var 'non_key'='0' with type 'int'")
}

func TestEnvIntWrong(t *testing.T) {
	expected, def := "string", 0

	t.Setenv(existingKey, expected)

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to int",
		func() { c.EnvInt(existingKey, def) },
	)
}

func TestEnvFloat(t *testing.T) {
	expected, def := 14.0, 0.0

	t.Setenv(existingKey, strconv.FormatFloat(expected, 'f', -1, 64))

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	exists := c.EnvFloat(existingKey, def)
	notExists := c.EnvFloat(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)

	assert.Contains(t, buf.String(), "loaded env var 'key'='14.00' with type 'float64'")
	assert.Contains(
		t,
		buf.String(),
		"loaded env var 'non_key'='0.00' with type 'float64'",
	)
}

func TestEnvFloatWrong(t *testing.T) {
	expected, def := "string", 0.0

	t.Setenv(existingKey, expected)

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to float64",
		func() { c.EnvFloat(existingKey, def) },
	)
}

func TestEnvBool(t *testing.T) {
	expected, def := true, false

	t.Setenv(existingKey, strconv.FormatBool(expected))

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	exists := c.EnvBool(existingKey, def)
	notExists := c.EnvBool(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)

	assert.Contains(t, buf.String(), "loaded env var 'key'='true' with type 'bool'")
	assert.Contains(
		t,
		buf.String(),
		"loaded env var 'non_key'='false' with type 'bool'",
	)
}

func TestEnvBoolWrong(t *testing.T) {
	expected, def := "string", false

	t.Setenv(existingKey, expected)

	var buf bytes.Buffer
	c := config.New(
		slog.New(
			logging.NewBufLogHandler(
				&buf,
				//nolint:exhaustruct //other fields are optional
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		),
	)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to bool",
		func() { c.EnvBool(existingKey, def) },
	)
}
