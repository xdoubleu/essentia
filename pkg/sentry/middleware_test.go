package sentry_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/config"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"
)

func testMiddleware(
	t *testing.T,
	middleware func(next http.Handler) http.Handler,
	req *http.Request,
	innerHandler func(w http.ResponseWriter, r *http.Request),
) *httptest.ResponseRecorder {
	t.Helper()

	var testResponse = []byte("bar")
	var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if innerHandler != nil {
			innerHandler(w, r)
		}
		_, _ = w.Write(testResponse)
		w.WriteHeader(http.StatusOK)
	})

	res := httptest.NewRecorder()
	middleware(testHandler).ServeHTTP(res, req)

	return res
}

func TestMiddleware(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	sentryMiddleware, err := sentrytools.Middleware(
		config.TestEnv,
		sentrytools.MockedSentryClientOptions(),
	)
	require.Nil(t, err)

	testMiddleware(
		t,
		sentryMiddleware,
		req,
		func(_ http.ResponseWriter, r *http.Request) {
			assert.NotNil(t, sentry.GetHubFromContext(r.Context()))
		},
	)
}
