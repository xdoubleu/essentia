package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/internal/mocks"
	"github.com/xdoubleu/essentia/pkg/context"
	"github.com/xdoubleu/essentia/pkg/middleware"
)

func testCORSHeaders(
	t *testing.T,
	handler func(next http.Handler) http.Handler,
	method string,
	allowedOrigins []string,
	exHeaders []string,
) {
	t.Helper()

	expectedMethods := []string{"GET", "POST", "PATCH", "DELETE"}

	req, _ := http.NewRequest(method, "http://example.com/foo", nil)
	req.Header["Origin"] = []string{"http://example.com"}

	if method == http.MethodOptions {
		req.Header["Access-Control-Request-Method"] = expectedMethods
		req.Header["Access-Control-Request-Headers"] = exHeaders[0:1]
	}

	res := testMiddleware(t, handler, req, nil)

	origins := res.Header()["Access-Control-Allow-Origin"]
	assert.Equal(t, allowedOrigins, origins)

	creds, _ := strconv.ParseBool(res.Header()["Access-Control-Allow-Credentials"][0])
	assert.True(t, creds)

	if method != http.MethodOptions {
		return
	}

	methods := res.Header()["Access-Control-Allow-Methods"]
	assert.Equal(t, expectedMethods, methods)

	headers := res.Header()["Access-Control-Allow-Headers"]
	assert.Equal(t, exHeaders[0:1], headers)
}

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

func TestCORS(t *testing.T) {
	allowedOrigins := []string{"http://example.com"}

	sentryHeaders := []string{"content-type", "baggage", "sentry-trace"}
	noSentryHeaders := []string{"content-type"}

	corsSentry := middleware.CORS(allowedOrigins, true)
	corsNoSentry := middleware.CORS(allowedOrigins, false)

	testCORSHeaders(t, corsSentry, http.MethodGet, allowedOrigins, sentryHeaders)
	testCORSHeaders(t, corsSentry, http.MethodOptions, allowedOrigins, sentryHeaders)

	testCORSHeaders(t, corsNoSentry, http.MethodGet, allowedOrigins, noSentryHeaders)
	testCORSHeaders(
		t,
		corsNoSentry,
		http.MethodOptions,
		allowedOrigins,
		noSentryHeaders,
	)
}

func TestErrors(t *testing.T) {
	showErrors := middleware.ShowErrors()

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	testMiddleware(
		t,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		},
		req,
		func(_ http.ResponseWriter, r *http.Request) {
			assert.False(t, context.ShowErrors(r.Context()))
		},
	)
	testMiddleware(t, showErrors, req, func(_ http.ResponseWriter, r *http.Request) {
		assert.True(t, context.ShowErrors(r.Context()))
	})
}

func TestLogger(t *testing.T) {
	mockedLogger := mocks.MockedLogger{}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	testMiddleware(
		t,
		middleware.Logger(mockedLogger.Logger()),
		req,
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	timeStr := time.Now().Format("2006-01-02")
	assert.Contains(t, mockedLogger.CapturedLogs(), timeStr)
	assert.Contains(
		t,
		mockedLogger.CapturedLogs(),
		fmt.Sprintf("status=%d", http.StatusOK),
	)
}

func TestRateLimit(t *testing.T) {
	bucketSize := 30
	rateLimiter := middleware.RateLimit(
		10,
		bucketSize,
		time.Second,
		500*time.Millisecond,
	)

	singleRequest := func() *httptest.ResponseRecorder {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
		req.RemoteAddr = "127.0.0.1:80"

		return testMiddleware(t, rateLimiter, req, nil)
	}

	for i := 0; i < bucketSize; i++ {
		res := singleRequest()
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	}
	time.Sleep(10 * time.Millisecond)

	res := singleRequest()
	assert.Equal(t, http.StatusTooManyRequests, res.Result().StatusCode)
	time.Sleep(2 * time.Second)

	res = singleRequest()
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}

func TestRecover(t *testing.T) {
	mockedLogger := mocks.MockedLogger{}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	res := testMiddleware(
		t,
		middleware.Recover(mockedLogger.Logger()),
		req,
		func(_ http.ResponseWriter, _ *http.Request) {
			panic("test")
		},
	)

	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, "close", res.Header()["Connection"][0])
	assert.Contains(t, mockedLogger.CapturedLogs(), "PANIC")
}
