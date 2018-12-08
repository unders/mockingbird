package mock

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// HTMLAdapter is used for tests
type HTMLAdapter struct {
	Code       int
	Body       []byte
	Err        error
	EnvErr     error
	ServiceErr error
}

//
// Business  Logic
//

// Dashboard returns the dashboard page
func (a HTMLAdapter) Dashboard() (code int, body []byte, err error) {
	b := append(a.Body, "dashboard page"...)
	return a.Code, b, a.Err
}

// Dashboard starts a test suite
func (a HTMLAdapter) RunTest() (id int, err error) {

	return a.Code, err
}

// ShowTestResult returns the ShowTestResult page
func (a HTMLAdapter) ShowTestResult(id int) (code int, body []byte, err error) {
	return a.Code, a.Body, a.Err
}

// ListTestResults returns the ListTestResults page
func (a HTMLAdapter) ListTestResults() (code int, body []byte, err error) {
	// FIXME: This should be removed, just testing the concept.
	result, err := sh.Output("mage", "test:all")

	status := sh.ExitStatus(err)

	fmt.Printf("result: %s\nStatus: %d", result, status)
	return a.Code, []byte(result), a.Err
}

// RunTestForService starts a test suite for a service
func (a HTMLAdapter) RunTestForService(service string) (id int, err error) {
	return a.Code, a.Err
}

//
// Validations
//
func (a HTMLAdapter) HasServiceError(env string) error {
	return a.ServiceErr
}

//
// Error pages
//

// ErrorNotFound returns error not found page
func (a HTMLAdapter) ErrorNotFound() (code int, body []byte) {
	b := append(a.Body, "Not Found"...)
	return a.Code, b
}

// ErrorClient returns an error page when client has made an error
func (a HTMLAdapter) ErrorClient(title, description string) (code int, body []byte) {
	b := append(a.Body, "Error Client"...)
	b = append(b, " title: "...)
	b = append(b, title...)
	b = append(b, " description: "...)
	b = append(b, description...)
	return a.Code, b
}

// ErrorServer returns an error page when server has an error
func (a HTMLAdapter) ErrorServer() (code int, body []byte) {
	b := append(a.Body, "Error Server"...)
	return a.Code, b
}
