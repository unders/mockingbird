package mock

import "github.com/unders/mockingbird/cmd/domain/mockinbird"

// HTMLAdapter is used for tests
type HTMLAdapter struct {
	Code       int
	Body       []byte
	Err        error
	EnvErr     error
	ServiceErr error
}

// Verifies that HTMLAdapter implements mockingbird.HTMLAdapter interface
var _ mockinbird.HTMLAdapter = HTMLAdapter{}

//
// Business  Logic
//

// Dashboard returns the dashboard page
func (a HTMLAdapter) Dashboard() (code int, body []byte, err error) {
	return a.Code, a.Body, a.Err
}

// Dashboard starts a test suite
func (a HTMLAdapter) RunTest(environment string) (id int, err error) {
	return a.Code, a.Err
}

// ShowTestResult returns the ShowTestResult page
func (a HTMLAdapter) ShowTestResult(environment string, id int) (code int, body []byte, err error) {
	return a.Code, a.Body, a.Err
}

// ListTestResults returns the ListTestResults page
func (a HTMLAdapter) ListTestResults(env string) (code int, body []byte, err error) {
	return a.Code, a.Body, a.Err
}

// RunTestForService starts a test suite for a service
func (a HTMLAdapter) RunTestForService(env, service string) (id int, err error) {
	return a.Code, a.Err
}

//
// Validations
//
func (a HTMLAdapter) HasEnvError(env string) error {
	return a.EnvErr
}

func (a HTMLAdapter) HasServiceError(env string) error {
	return a.ServiceErr
}

//
// Error pages
//

// ErrorNotFound returns error not found page
func (a HTMLAdapter) ErrorNotFound() (code int, body []byte) {
	return a.Code, a.Body
}

// ErrorClient returns an error page when client has made an error
func (a HTMLAdapter) ErrorClient(title, description string) (code int, body []byte) {
	return a.Code, a.Body
}

// ErrorServer returns an error page when server has an error
func (a HTMLAdapter) ErrorServer() (code int, body []byte) {
	return a.Code, a.Body
}
