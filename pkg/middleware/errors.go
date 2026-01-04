package middleware

import (
	"net/http"

	"github.com/xdoubleu/essentia/internal/helpers"
	"github.com/xdoubleu/essentia/pkg/contexttools"
)

// ShowErrors is middleware used to show errors.
// When used errors handled by [httptools.ServerErrorResponse] will be shown.
// Otherwise these will be hidden.
func ShowErrors() helpers.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(contexttools.WithShownErrors(r.Context()))
			next.ServeHTTP(w, r)
		})
	}
}
