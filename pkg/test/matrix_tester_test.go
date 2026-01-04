package test_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/xdoubleu/essentia/pkg/communication/httptools"
	"github.com/xdoubleu/essentia/pkg/errortools"
	"github.com/xdoubleu/essentia/pkg/parse"
	"github.com/xdoubleu/essentia/pkg/test"
)

func matrixTestHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cookie")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	if cookie != nil && cookie.Value == "value" {
		http.SetCookie(w, &http.Cookie{Name: "cookie2", Value: cookie.Value})
		httptools.UnauthorizedResponse(
			w,
			r,
			errortools.NewUnauthorizedError(errors.New("unauthorized")),
		)
		return
	}

	param, err := parse.QueryParam(r, "param", "default", nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	if param == "value" {
		httptools.ForbiddenResponse(w, r)
		return
	}

	var data map[string]string
	err = httptools.ReadJSON(r.Body, &data)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	if data["test"] == "error" {
		httptools.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			map[string]any{"message": "test"},
		)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func TestMatrixTester(t *testing.T) {
	baseRequest := test.CreateRequestTester(
		http.HandlerFunc(matrixTestHandler),
		http.MethodGet,
		"",
	)
	mt := test.CreateMatrixTester()

	tReq1 := baseRequest.Copy()
	tReq1.SetData(map[string]string{
		"test": "error",
	})

	tRes1 := test.NewCaseResponse(http.StatusBadRequest, nil, errortools.NewErrorDto(
		http.StatusBadRequest,
		map[string]any{"message": "test"},
	))

	mt.AddTestCase(tReq1, tRes1)

	tReq2 := baseRequest.Copy()
	tReq2.SetQuery(url.Values{
		"param": {"value"},
	})

	mt.AddTestCase(tReq2, test.NewCaseResponse(http.StatusForbidden, nil, nil))

	tReq3 := baseRequest.Copy()
	tReq3.AddCookie(&http.Cookie{Name: "cookie", Value: "value"})

	tRes3 := test.NewCaseResponse(
		http.StatusUnauthorized,
		[]*http.Cookie{{Name: "cookie2", Value: "value"}},
		nil,
	)

	mt.AddTestCase(tReq3, tRes3)

	mt.Do(t)
}
