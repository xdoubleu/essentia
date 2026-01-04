package parse_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/parse"
)

func TestURLParamOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.SetPathValue("pathValue", "test")

	result, err := parse.URLParam[string](req, "pathValue", nil)

	assert.Equal(t, "test", result)
	assert.Equal(t, nil, err)
}

func TestURLParamNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	result, err := parse.URLParam[string](req, "pathValue", nil)

	assert.Equal(t, "", result)
	assert.Equal(t, errors.New("missing URL param 'pathValue'"), err)
}

func TestRequiredQueryParamOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.URL.RawQuery = "queryParam=test"

	result, err := parse.RequiredQueryParam[string](req, "queryParam", nil)

	assert.Equal(t, "test", result)
	assert.Equal(t, nil, err)
}

func TestRequiredQueryParamNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	result, err := parse.RequiredQueryParam[string](req, "queryParam", nil)

	assert.Equal(t, "", result)
	assert.Equal(t, errors.New("missing query param 'queryParam'"), err)
}

func TestQueryParamOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.URL.RawQuery = "queryParam=test"

	result, err := parse.QueryParam[string](req, "queryParam", "default", nil)

	assert.Equal(t, "test", result)
	assert.Equal(t, nil, err)
}

func TestQueryParamNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	result, err := parse.QueryParam[string](req, "queryParam", "default", nil)

	assert.Equal(t, "default", result)
	assert.Equal(t, nil, err)
}

func TestRequiredArrayQueryParamOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.URL.RawQuery = "queryParam=test1,test2"

	result, err := parse.RequiredArrayQueryParam[string](req, "queryParam", nil)

	assert.Equal(t, []string{"test1", "test2"}, result)
	assert.Equal(t, nil, err)
}

func TestRequiredArrayQueryParamNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	result, err := parse.RequiredArrayQueryParam[string](req, "queryParam", nil)

	assert.Equal(t, []string{}, result)
	assert.Equal(t, errors.New("missing query param 'queryParam'"), err)
}

func TestArrayQueryParamOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.URL.RawQuery = "queryParam=test1,test2"

	result, err := parse.ArrayQueryParam(
		req,
		"queryParam",
		[]string{"default1", "default2"},
		nil,
	)

	assert.Equal(t, []string{"test1", "test2"}, result)
	assert.Equal(t, nil, err)
}

func TestArrayQueryParamNOK(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)

	result, err := parse.ArrayQueryParam(
		req,
		"queryParam",
		[]string{"default1", "default2"},
		nil,
	)

	assert.Equal(t, []string{"default1", "default2"}, result)
	assert.Equal(t, nil, err)
}

// also covers parse func behavior of others.
func TestArrayQueryParamFailedParseFunc(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.URL.RawQuery = "queryParam=test1,test2"

	result, err := parse.ArrayQueryParam(
		req,
		"queryParam",
		[]int{1, 2},
		parse.Int(false, true),
	)

	assert.Equal(t, []int{}, result)
	assert.Equal(
		t,
		errors.New(
			"invalid query param 'queryParam' with value 'test1', should be an integer",
		),
		err,
	)
}
