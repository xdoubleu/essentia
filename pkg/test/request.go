package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/xdoubleu/essentia/v2/pkg/communication/httptools"
)

// ContentType is used to set the "Content-Type" header of requests.
// This also has an influence on the encoding of your data.
type ContentType = string

const (
	// JSONContentType sets the "Content-Type" header to "application/json".
	JSONContentType ContentType = "application/json"
	// FormContentType sets the "Content-Type" header to
	// "application/x-www-form-urlencoded".
	FormContentType = "application/x-www-form-urlencoded"
)

// A RequestTester is used to test a certain HTTP request.
type RequestTester struct {
	handler         http.Handler
	ts              *httptest.Server
	contentType     ContentType
	method          string
	path            string
	data            any
	query           url.Values
	cookies         []*http.Cookie
	followRedirects bool
	RawRequest      *http.Request
}

// CreateRequestTester creates a new [RequestTester].
func CreateRequestTester(
	handler http.Handler,
	method, path string,
	pathValues ...any,
) RequestTester {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(pathValues) > 0 {
		path = fmt.Sprintf(path, pathValues...)
	}

	return RequestTester{
		handler,
		nil,
		JSONContentType,
		method,
		path,
		nil,
		make(url.Values),
		[]*http.Cookie{},
		true,
		nil,
	}
}

// SetTestServer sets the test server of a [RequestTester].
// This allows to reuse existing test servers.
func (tReq *RequestTester) SetTestServer(ts *httptest.Server) {
	tReq.ts = ts
}

// SetContentType sets the content type of a [RequestTester].
func (tReq *RequestTester) SetContentType(contentType ContentType) {
	tReq.contentType = contentType
}

// SetData sets the request data of a [RequestTester].
func (tReq *RequestTester) SetData(data any) {
	tReq.data = data
}

// SetQuery sets the query of a [RequestTester].
func (tReq *RequestTester) SetQuery(query url.Values) {
	tReq.query = query
}

// AddCookie adds a cookie to a [RequestTester].
// Can be used multiple times for adding several cookies.
func (tReq *RequestTester) AddCookie(cookie *http.Cookie) {
	tReq.cookies = append(tReq.cookies, cookie)
}

// SetFollowRedirect configures the request to follow or ignore redirects.
func (tReq *RequestTester) SetFollowRedirect(follow bool) {
	tReq.followRedirects = follow
}

// Do executes a [RequestTester] returning the response of a request
// and providing the returned data to rsData.
func (tReq *RequestTester) Do(t *testing.T) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	var err error
	var ts *httptest.Server

	if tReq.ts != nil {
		ts = tReq.ts
	} else if tReq.handler != nil {
		ts = httptest.NewServer(tReq.handler)
		defer ts.Close()
	}

	if ts == nil {
		panic("handler nor test server has been set")
	}

	if tReq.data != nil {
		switch tReq.contentType {
		case JSONContentType:
			var body []byte
			body, err = json.Marshal(tReq.data)
			bodyReader = bytes.NewReader(body)
		case FormContentType:
			var query url.Values
			query, err = httptools.WriteForm(tReq.data)
			bodyReader = strings.NewReader(query.Encode())
		}

		if err != nil {
			t.Errorf("error when marshalling body: %v", err)
			t.FailNow()
			return nil
		}
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		tReq.method,
		fmt.Sprintf("%s/%s", ts.URL, tReq.path),
		bodyReader,
	)

	if err != nil {
		t.Errorf("error when creating request: %v", err)
		t.FailNow()
		return nil
	}

	if len(tReq.query) > 0 {
		query := req.URL.Query()

		for key, values := range tReq.query {
			for _, value := range values {
				query.Add(key, value)
			}
		}

		req.URL.RawQuery = query.Encode()
	}

	for _, cookie := range tReq.cookies {
		req.AddCookie(cookie)
	}

	req.Header.Set("Content-Type", tReq.contentType)

	tReq.RawRequest = req

	client := ts.Client()

	if !tReq.followRedirects {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	rs, err := client.Do(req)
	if err != nil {
		t.Errorf("error when making request: %v", err)
		t.FailNow()
		return nil
	}

	return rs
}

// Copy creates a copy of a [RequestTester] in order to easily test similar requests.
func (tReq RequestTester) Copy() RequestTester {
	return RequestTester{
		handler:         tReq.handler,
		ts:              tReq.ts,
		contentType:     tReq.contentType,
		method:          tReq.method,
		path:            tReq.path,
		data:            tReq.data,
		query:           tReq.query,
		cookies:         tReq.cookies,
		followRedirects: tReq.followRedirects,
		RawRequest:      tReq.RawRequest,
	}
}
