package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
)

func TestCSV(t *testing.T) {
	res := httptest.NewRecorder()

	headers := []string{"h1", "h2"}
	data := [][]string{
		{
			"c11",
			"c12",
		},
	}
	expectedOutput := [][]string{}
	expectedOutput = append(expectedOutput, headers)
	expectedOutput = append(expectedOutput, data...)

	writeCSV := func(w http.ResponseWriter, _ *http.Request) {
		err := httptools.WriteCSV(w, "test", headers, data)
		require.Nil(t, err)
	}
	http.HandlerFunc(writeCSV).ServeHTTP(res, nil)

	records, err := httptools.ReadCSV(res.Body)

	require.Nil(t, err)
	assert.Equal(t, expectedOutput, records)
}

func TestJSON(t *testing.T) {
	res := httptest.NewRecorder()

	data := [][]string{
		{
			"h1",
			"h2",
		},
		{
			"c11",
			"c12",
		},
	}

	writeJSON := func(w http.ResponseWriter, _ *http.Request) {
		err := httptools.WriteJSON(w, http.StatusOK, data, nil)
		require.Nil(t, err)
	}
	http.HandlerFunc(writeJSON).ServeHTTP(res, nil)

	var result [][]string
	err := httptools.ReadJSON(res.Body, &result)

	require.Nil(t, err)
	assert.Equal(t, data, result)
}
