package html

import "github.com/unders/mockingbird/server/domain/mockingbird"

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
	App mockingbird.App
}

// Verifies that *Adapter implements mockingbird.HTMLAdapter interface
var _ mockingbird.HTMLAdapter = &Adapter{}

//
// Business  Logic
//

// Dashboard returns the dashboard page
func (a Adapter) Dashboard() (code int, body []byte, err error) {
	return 200, []byte("Dashboard"), nil
}

// RunTest starts a test suite
func (a Adapter) RunTest() (id string, code int, cody []byte, err error) {
	return "100", 200, []byte("RunTest"), nil
}

// ShowTestResult returns the ShowTestResult page
func (a Adapter) ShowTestResult(id string) (code int, body []byte, err error) {
	return 200, []byte("Test result page"), nil
}

// ListTestResults returns the ListTestResults page
func (a Adapter) ListTestResults(service string) (code int, body []byte, err error) {
	return 200, []byte("List test result"), nil
}

// RunTestForService starts a test suite for a service
func (a Adapter) RunTestForService(service string) (id string, code int, body []byte, err error) {
	return "50045", 200, []byte("List test result"), nil
}

//
// Error pages
//

// ErrorNotFound returns error not found page
func (a Adapter) ErrorNotFound() (body []byte) {
	return []byte("ErrorNotFound page")
}

// InvalidURL returns error an Invalid URL page
func (a Adapter) InvalidURL() (body []byte) {
	return []byte("InvalidURL page")
}
