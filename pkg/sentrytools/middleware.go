package sentrytools

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/xdoubleu/essentia/v3/internal/helpers"
	"github.com/xdoubleu/essentia/v3/pkg/config"
)

// Init initializes Sentry and returns a hub clone suitable for use on
// background contexts (e.g. job queues). Returns nil, nil when the DSN is
// empty or env is [config.TestEnv] — in those cases no initialization is
// performed and Middleware() handles its own setup.
// Must be called before Middleware().
func Init(env string, options sentry.ClientOptions) (*sentry.Hub, error) {
	if env == config.TestEnv || options.Dsn == "" {
		return nil, nil //nolint:nilnil //Sentry disabled is not an error
	}

	if err := sentry.Init(options); err != nil {
		return nil, err
	}

	return sentry.CurrentHub().Clone(), nil
}

// Middleware is middleware used to configure and enable Sentry.
// Call [Init] at application startup before using this middleware.
// When env is [config.TestEnv], a mocked [sentry.Hub] will be used and
// Sentry is self-initialized with mock options.
func Middleware(env string) (helpers.Middleware, error) {
	isTestEnv := env == config.TestEnv

	if isTestEnv {
		if err := sentry.Init(MockedSentryClientOptions()); err != nil {
			return nil, err
		}
	}

	//nolint:exhaustruct //other fields are optional
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	if isTestEnv {
		return func(next http.Handler) http.Handler {
			return sentryHandler.Handle(useMockedHub(next))
		}, nil
	}

	return func(next http.Handler) http.Handler {
		return sentryHandler.Handle(next)
	}, nil
}

func useMockedHub(next http.Handler) http.Handler {
	mockedHub := MockedSentryHub()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sentry.SetHubOnContext(r.Context(), mockedHub)
		next.ServeHTTP(w, r)
	})
}
