package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/xdoubleu/essentia/internal/shared"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/context"
)

// Logger is middleware used to add a logger to
// the context and log every request and their duration.
func Logger(logger *slog.Logger) shared.Middleware {
	return func(next http.Handler) http.Handler {
		return loggerHandler(logger, next)
	}
}

func loggerHandler(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httptools.NewResponseWriter(w)
		t := time.Now()

		r = r.WithContext(context.WithLogger(r.Context(), logger))

		next.ServeHTTP(rw, r)

		logger.Info(
			"processed request",
			slog.Int("status", rw.Status()),
			slog.String("endpoint", r.RequestURI),
			slog.Duration("duration", time.Since(t)),
		)
	})
}
