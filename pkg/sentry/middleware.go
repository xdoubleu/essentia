package sentry

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/xdoubleu/essentia/internal/shared"
	"github.com/xdoubleu/essentia/pkg/config"
)

// Middleware is middleware used to configure and enable Sentry.
// When env is [config.TestEnv], a mocked [sentry.Hub] will be used.
func Middleware(
	env string,
	clientOptions sentry.ClientOptions,
) (shared.Middleware, error) {
	isTestEnv := env == config.TestEnv

	if isTestEnv {
		clientOptions = MockedSentryClientOptions()
	}

	sentryHandler, err := getSentryHandler(clientOptions)
	if err != nil {
		return nil, err
	}

	if isTestEnv {
		return func(next http.Handler) http.Handler {
			return sentryHandler.Handle(useMockedHub(next))
		}, nil
	}

	return func(next http.Handler) http.Handler {
		return sentryHandler.Handle(next)
	}, nil
}

func getSentryHandler(clientOptions sentry.ClientOptions) (*sentryhttp.Handler, error) {
	err := sentry.Init(clientOptions)

	if err != nil {
		return nil, err
	}

	//nolint:exhaustruct //other fields are optional
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}), nil
}

func useMockedHub(next http.Handler) http.Handler {
	mockedHub := MockedSentryHub()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sentry.SetHubOnContext(r.Context(), mockedHub)
		next.ServeHTTP(w, r)
	})
}
