package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
)

func TestRedirect(t *testing.T) {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			httptools.RedirectWithError(w, r, "url", errors.New("test"))
		},
	)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "", nil)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusSeeOther, res.Result().StatusCode)
}
