package test_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/communication/httptools"
	"github.com/xdoubleu/essentia/pkg/test"
)

type Test struct {
	Key string
}

func testHandlerJSON(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cookiename")
	if err != nil || cookie.Value != "value" {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	var data Test
	err = httptools.ReadJSON(r.Body, &data)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func TestRequestTesterJson(t *testing.T) {
	reqData := Test{
		Key: "data",
	}

	tReq := test.CreateRequestTester(
		http.HandlerFunc(testHandlerJSON),
		http.MethodPost,
		"/test/%d",
		1,
	)
	tReq.AddCookie(&http.Cookie{Name: "cookiename", Value: "value"})
	tReq.SetData(reqData)

	rs := tReq.Do(t)

	var rsData Test
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, reqData, rsData)
}

func testHandlerForm(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cookiename")
	if err != nil || cookie.Value != "value" {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	var data Test
	err = httptools.ReadForm(r, &data)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func TestRequestTesterForm(t *testing.T) {
	reqData := Test{
		Key: "data",
	}

	tReq := test.CreateRequestTester(
		http.HandlerFunc(testHandlerForm),
		http.MethodPost,
		"/test/%d",
		1,
	)
	tReq.AddCookie(&http.Cookie{Name: "cookiename", Value: "value"})
	tReq.SetContentType(test.FormContentType)
	tReq.SetData(reqData)

	rs := tReq.Do(t)

	var rsData Test
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, reqData, rsData)
}

func TestRequestTesterTestServer(t *testing.T) {
	reqData := Test{
		Key: "data",
	}

	ts := httptest.NewServer(http.HandlerFunc(testHandlerJSON))

	tReq := test.CreateRequestTester(
		nil,
		http.MethodGet,
		"/test/%d",
		1,
	)
	tReq.AddCookie(&http.Cookie{Name: "cookiename", Value: "value"})
	tReq.SetData(reqData)
	tReq.SetTestServer(ts)

	rs := tReq.Do(t)

	var rsData Test
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, reqData, rsData)
}

func TestRequestTesterNoTestServerOrHandler(t *testing.T) {
	tReq := test.CreateRequestTester(nil, http.MethodGet, "")

	assert.PanicsWithValue(
		t,
		"handler nor test server has been set",
		func() { tReq.Do(t) },
	)
}
