package httptools_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/communication/httptools"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/contexttools"
	"github.com/xdoubleu/essentia/pkg/errortools"
	"github.com/xdoubleu/essentia/pkg/sentrytools"
)

func testError(t *testing.T, handler http.HandlerFunc) (int, errortools.ErrorDto) {
	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	return testErrorWithReq(t, handler, req)
}

func testErrorWithReq(
	t *testing.T,
	handler http.HandlerFunc,
	req *http.Request,
) (int, errortools.ErrorDto) {
	t.Helper()

	res := httptest.NewRecorder()

	sentryMiddleware, err := sentrytools.Middleware(
		config.TestEnv,
		sentrytools.MockedSentryClientOptions(),
	)
	require.Nil(t, err)

	sentryMiddleware(handler).ServeHTTP(res, req)

	var errorDto errortools.ErrorDto
	err = httptools.ReadJSON(res.Result().Body, &errorDto)
	require.Nil(t, err)

	return res.Result().StatusCode, errorDto
}

func TestServerErrorResponseObfuscated(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.HandleError(w, r, errors.New("test"))
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, errortools.MessageInternalServerError, errorDto.Message)
}

func TestServerErrorResponseShown(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.HandleError(w, r, errors.New("test"))
	}

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	req = req.WithContext(contexttools.WithShownErrors(req.Context()))

	statusCode, errorDto := testErrorWithReq(t, handler, req)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "test", errorDto.Message)
}

func TestBadRequestResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.HandleError(
			w,
			r,
			errortools.NewBadRequestError(errors.New("test")),
		)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "test", errorDto.Message)
}

func TestRateLimitExceededResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.RateLimitExceededResponse(w, r)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusTooManyRequests, statusCode)
	assert.Equal(t, errortools.MessageTooManyRequests, errorDto.Message)
}

func TestUnauthorizedResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.HandleError(
			w,
			r,
			errortools.NewUnauthorizedError(errors.New("test")),
		)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "test", errorDto.Message)
}

func TestForbiddenResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.ForbiddenResponse(w, r)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusForbidden, statusCode)
	assert.Equal(t, errortools.MessageForbidden, errorDto.Message)
}

func TestConflictResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.HandleError(
			w,
			r,
			errortools.NewConflictError("resource", "value", "field"),
		)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, map[string]any{
		"field": "resource with field 'value' already exists",
	}, errorDto.Message)
}

func TestNotFoundResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.HandleError(
			w,
			r,
			errortools.NewNotFoundError("resource", "value", "field"),
		)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, map[string]any{
		"field": "resource with field 'value' doesn't exist",
	}, errorDto.Message)
}

func TestFailedValidationResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.FailedValidationResponse(w, r, map[string]string{
			"field": "invalid value",
		})
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, map[string]any{
		"field": "invalid value",
	}, errorDto.Message)
}
