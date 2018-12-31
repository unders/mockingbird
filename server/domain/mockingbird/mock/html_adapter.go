package mock

import (
	"net/http"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/pkg/errs"
)

// HTMLAdapter is used for tests
type HTMLAdapter struct {
	Code       int
	Body       []byte
	Err        error
	IDErr      error
	ServiceErr error
}

// Verifies that HTMLAdapter implements mockingbird.HTMLAdapter interface
var _ mockingbird.HTMLAdapter = HTMLAdapter{}

//
// Business  Logic
//

// Dashboard returns the dashboard page
func (a HTMLAdapter) Dashboard() (code int, body []byte, err error) {
	b := append(a.Body, "dashboard page"...)
	if a.Err != nil {
		b = append(b, " with error"...)
	}
	return a.Code, b, a.Err
}

// ListTests returns the test results page
func (a HTMLAdapter) ListTests(pageToken string) (code int, body []byte, err error) {
	b := append(a.Body, "list test result page"...)

	if a.Err != nil {
		b = append(b, " with error"...)
	}
	return a.Code, b, a.Err
}

// ShowTest returns the test result page
func (a HTMLAdapter) ShowTest(id mockingbird.ULID) (code int, body []byte, err error) {
	if a.IDErr != nil {
		b := append(a.Body, "HasIdError for "...)
		b = append(b, id...)
		return a.Code, b, a.IDErr
	}

	b := append(a.Body, "test result page for id="...)
	b = append(b, id...)
	if a.Err != nil {
		b = append(b, " with error"...)
	}
	return a.Code, b, a.Err
}

// RunTest starts a test suite
func (a HTMLAdapter) RunTest(ts mockingbird.TestSuite) (id mockingbird.ULID, code int, cody []byte, err error) {
	if ts == "" {
		return "", http.StatusNotFound, []byte("RunTest failed with 404 not found"), errs.NotFound("Not found")
	}

	b := a.Body
	if a.Err != nil {
		b = append(b, "with runTest error"...)
	}
	return "test-suite-id", a.Code, b, a.Err
}

//
// Error pages
//

// ErrorNotFound returns error not found page
func (a HTMLAdapter) ErrorNotFound() (body []byte) {
	b := append(a.Body, "Not Found"...)
	return b
}

// InvalidURL returns error an Invalid URL page
func (a HTMLAdapter) InvalidURL() (body []byte) {
	b := append(a.Body, "Invalid URL page"...)
	return b
}
