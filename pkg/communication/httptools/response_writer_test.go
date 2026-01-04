package httptools_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/communication/httptools"
)

func TestStatus(t *testing.T) {
	res := httptest.NewRecorder()
	rw := httptools.NewResponseWriter(res)
	rw.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, rw.Status())
}
