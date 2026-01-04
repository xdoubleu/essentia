// Package config provides functions which can be used to
// extract environment variables and parse them to the right type.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xdoubleu/essentia/internal/shared"
)

// Parser parses the config provided through environment variables.
type Parser struct {
	logger *slog.Logger
}

const (
	// ProdEnv can be used as value when reading out the type of environment.
	ProdEnv string = "production"
	// TestEnv can be used as value when reading out the type of environment.
	TestEnv string = "test"
	// DevEnv can be used as value when reading out the type of environment.
	DevEnv string = "development"
)

const errorMessage = "can't convert env var '%s' with value '%s' to %s"

// New returns a new Parser and loads environment variables that
// could be provided using a .env file (particularly useful during development).
func New(logger *slog.Logger) Parser {
	_ = godotenv.Load()

	return Parser{
		logger: logger,
	}
}

func (c Parser) baseEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return ""
	}

	return value
}

func (c Parser) logValue(valType string, key string, value any) {
	strVal, err := shared.AnyToString(value)
	if err != nil {
		panic(err)
	}

	c.logger.Info(
		fmt.Sprintf("loaded env var '%s'='%s' with type '%s'", key, strVal, valType),
	)
}

// EnvStr extracts a string environment variable.
func (c Parser) EnvStr(key string, defaultValue string) string {
	value := c.baseEnv(key)
	if len(value) == 0 {
		value = defaultValue
	}

	c.logValue("string", key, value)
	return value
}

// EnvStrArray extracts a string
// array environment variable. The values should be separated by ','.
func (c Parser) EnvStrArray(key string, defaultValue []string) []string {
	value := defaultValue

	strVal := c.baseEnv(key)
	if len(strVal) != 0 {
		value = strings.Split(strVal, ",")
	}

	c.logValue("string array", key, value)
	return value
}

// EnvInt extracts an integer environment variable.
func (c Parser) EnvInt(key string, defaultValue int) int {
	value := defaultValue

	strVal := c.baseEnv(key)
	if len(strVal) != 0 {
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			panic(fmt.Sprintf(errorMessage, key, strVal, "int"))
		}
		value = intVal
	}

	c.logValue("int", key, value)
	return value
}

// EnvFloat extracts a float environment variable.
func (c Parser) EnvFloat(key string, defaultValue float64) float64 {
	value := defaultValue

	strVal := c.baseEnv(key)
	if len(strVal) != 0 {
		floatVal, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			panic(fmt.Sprintf(errorMessage, key, strVal, "float64"))
		}
		value = floatVal
	}

	c.logValue("float64", key, value)
	return value
}

// EnvBool extracts a boolean environment variable.
func (c Parser) EnvBool(key string, defaultValue bool) bool {
	value := defaultValue

	strVal := c.baseEnv(key)
	if len(strVal) != 0 {
		boolVal, err := strconv.ParseBool(strVal)
		if err != nil {
			panic(fmt.Sprintf(errorMessage, key, strVal, "bool"))
		}
		value = boolVal
	}

	c.logValue("bool", key, value)
	return value
}
