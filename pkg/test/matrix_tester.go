package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
)

// CaseResponse is used to compare to the actual response
// of a [RequestTester] when used by a [MatrixTester].
type CaseResponse struct {
	statusCode int
	cookies    []*http.Cookie
	body       *map[string]any
}

// NewCaseResponse returns a new [CaseResponse].
func NewCaseResponse(statusCode int, cookies []*http.Cookie, body any) CaseResponse {
	caseResponse := CaseResponse{
		statusCode: statusCode,
		cookies:    nil,
		body:       nil,
	}

	if cookies != nil {
		caseResponse.cookies = cookies
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(data, &caseResponse.body)
		if err != nil {
			panic(err)
		}
	}

	return caseResponse
}

// MatrixTester is used for executing matrix tests.
type MatrixTester struct {
	testCases map[*RequestTester]CaseResponse
}

// CreateMatrixTester creates a new [MatrixTester].
func CreateMatrixTester() MatrixTester {
	return MatrixTester{
		testCases: make(map[*RequestTester]CaseResponse),
	}
}

// AddTestCase adds a test case which consists of
// a [RequestTester] and a [CaseResponse].
// When executing [Do] the [RequestTester] will be executed
// and its response will be compared to the provided [CaseResponse].
func (mt *MatrixTester) AddTestCase(
	tReq RequestTester,
	tRes CaseResponse) {
	mt.testCases[&tReq] = tRes
}

// Do executes a [MatrixTester].
func (mt MatrixTester) Do(t *testing.T) {
	t.Helper()

	for tReq, tRes := range mt.testCases {
		var rs *http.Response
		var rsData map[string]any

		if tRes.body == nil {
			rs = tReq.Do(t)
			defer rs.Body.Close()
		} else {
			rs = tReq.Do(t)
			err := httptools.ReadJSON(rs.Body, &rsData)
			require.Nil(t, err)

			defer rs.Body.Close()
		}

		assert.Equal(t, tRes.statusCode, rs.StatusCode)

		if tRes.cookies != nil {
			for _, cookie := range tRes.cookies {
				found := false

				for _, acCookie := range rs.Cookies() {
					found = cookie.String() == acCookie.String()
				}

				assert.True(t, found)
			}
		}

		if rsData != nil {
			assert.Equal(t, *tRes.body, rsData)
		}
	}
}
