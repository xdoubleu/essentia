package middleware

import (
	"github.com/rs/cors"
	"github.com/xdoubleu/essentia/internal/shared"
)

// CORS is middleware used to apply CORS settings.
func CORS(allowedOrigins []string, useSentry bool) shared.Middleware {
	allowedHeaders := []string{"content-type"}
	if useSentry {
		allowedHeaders = append(allowedHeaders, "baggage", "sentry-trace")
	}

	//nolint:exhaustruct //other fields are optional
	cors := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   allowedHeaders,
	})

	return cors.Handler
}
