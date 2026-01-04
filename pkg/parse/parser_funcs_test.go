package parse_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/pkg/parse"
)

func TestURLParamString(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", "value")

	result, err := parse.URLParam(req, "pathValue", parse.String)

	assert.Equal(t, "value", result)
	assert.Equal(t, nil, err)
}

func TestURLParamUUIDOK(t *testing.T) {
	val, _ := uuid.NewV7()

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", val.String())

	result, err := parse.URLParam(req, "pathValue", parse.UUID)

	assert.Equal(t, val.String(), result)
	assert.Equal(t, nil, err)
}

func TestURLParamUUIDNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", "notuuid")

	result, err := parse.URLParam(req, "pathValue", parse.UUID)

	assert.Equal(t, "", result)
	assert.Equal(
		t,
		errors.New(
			"invalid URL param 'pathValue' with value 'notuuid', should be a UUID",
		),
		err,
	)
}

func TestURLParamInt64OK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", fmt.Sprintf("%d", 1))

	result, err := parse.URLParam(req, "pathValue", parse.Int64(false, true))

	assert.Equal(t, int64(1), result)
	assert.Equal(t, nil, err)
}

func TestURLParamIntOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", fmt.Sprintf("%d", 1))

	result, err := parse.URLParam(req, "pathValue", parse.Int(false, true))

	assert.Equal(t, 1, result)
	assert.Equal(t, nil, err)
}

func TestURLParamIntNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", "notint")

	result, err := parse.URLParam(req, "pathValue", parse.Int(false, true))

	assert.Equal(t, 0, result)
	assert.Equal(
		t,
		errors.New(
			"invalid URL param 'pathValue' with value 'notint', should be an integer",
		),
		err,
	)
}

func TestURLParamIntLTZero(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", fmt.Sprintf("%d", -1))

	result, err := parse.URLParam(req, "pathValue", parse.Int(true, true))

	assert.Equal(t, 0, result)
	assert.Equal(
		t,
		errors.New(
			"invalid URL param 'pathValue' with value '-1', can't be less than '0'",
		),
		err,
	)
}

func TestURLParamIntZero(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", fmt.Sprintf("%d", 0))

	result, err := parse.URLParam(req, "pathValue", parse.Int(true, false))

	assert.Equal(t, 0, result)
	assert.Equal(
		t,
		errors.New("invalid URL param 'pathValue' with value '0', can't be '0'"),
		err,
	)
}

func TestURLParamDateOK(t *testing.T) {
	datetime := time.Now().Format("2006-01-02")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", datetime)

	result, err := parse.URLParam(req, "pathValue", parse.Date("2006-01-02"))

	expected, _ := time.Parse("2006-01-02", datetime)
	assert.Equal(t, expected, result)
	assert.Equal(t, nil, err)
}

func TestURLParamDateNOK(t *testing.T) {
	datetime := time.Now().Format("01-02-2006")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", datetime)

	result, err := parse.URLParam(req, "pathValue", parse.Date("2006-01-02"))

	expected, _ := time.Parse("2006-01-02", datetime)
	assert.Equal(t, expected, result)
	assert.Equal(
		t,
		fmt.Errorf(
			"invalid URL param 'pathValue' with value '%s', need format '2006-01-02'",
			datetime,
		),
		err,
	)
}
