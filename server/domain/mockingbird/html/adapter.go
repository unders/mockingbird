package html

import (
	"net/http"

	"github.com/unders/mockingbird/server/pkg/errs"

	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// Adapter implements the mockingbird.HTMLAdapter interface
//
//
// Note:
//
//         * It adapts between the http.Handler and the mockingbird.App
//         * This adapter returns HTML pages.
//
//
type Adapter struct {
	Tmpl *Template
	App  mockingbird.App
}

// Verifies that *Adapter implements mockingbird.HTMLAdapter interface
var _ mockingbird.HTMLAdapter = &Adapter{}

//
// Business  Logic
//

// Dashboard returns the dashboard page
func (a Adapter) Dashboard() (code int, body []byte, err error) {
	d, err := a.App.Dashboard()
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}

	b, err := a.Tmpl.Dashboard(d)
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}

	return 200, b, nil
}

// ListTest returns the ListTest page
func (a Adapter) ListTests(pageToken string) (code int, body []byte, err error) {
	tests, err := a.App.ListTests(pageToken)
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}
	b, err := a.Tmpl.ListTest(tests)
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}

	return 200, b, nil
}

// ShowTest returns the ShowTest page
func (a Adapter) ShowTest(id mockingbird.ULID) (code int, body []byte, err error) {
	test, err := a.App.ShowTest(id)
	if errs.IsNotFound(err) {
		return http.StatusNotFound, a.Tmpl.ErrorNotFound(), err
	}
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}
	b, err := a.Tmpl.ShowTest(test)
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}

	return 200, b, nil
}

// ShowTestSuites returns the ShowTestSuites page
func (a Adapter) ShowTestSuites() (code int, body []byte, err error) {
	ts := a.App.ShowTestSuites()
	b, err := a.Tmpl.ShowTestSuites(ts)
	if err != nil {
		return http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}

	return 200, b, nil
}

// RunTest starts a test suite
func (a Adapter) RunTest(ts mockingbird.TestSuite) (id mockingbird.ULID, code int, body []byte, err error) {
	ulid, err := a.App.RunTest(ts)
	if errs.IsNotFound(err) {
		return ulid, http.StatusNotFound, a.Tmpl.ErrorNotFound(), err
	}
	if err != nil {
		return ulid, http.StatusInternalServerError, a.Tmpl.InternalError(), err
	}

	// When no error, code, body, nil are ignored...
	return ulid, 200, nil, nil
}

//
// Error pages
//

// ErrorNotFound returns error not found page
func (a Adapter) ErrorNotFound() (body []byte) {
	return a.Tmpl.ErrorNotFound()
}

// InvalidURL returns error an Invalid URL page
func (a Adapter) InvalidURL() (body []byte) {
	return a.Tmpl.InvalidURL()
}
