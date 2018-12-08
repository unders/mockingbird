package mockingbird

import "github.com/unders/mockingbird/server/domain/mockingbird/mock"

// HTMLAdapter interface for working with HTML pages
type HTMLAdapter interface {
	//
	// Business  Logic
	//
	Dashboard() (code int, body []byte, err error)
	RunTest() (id int, err error)
	ShowTestResult(id int) (code int, body []byte, err error)
	ListTestResults() (code int, body []byte, err error)
	RunTestForService(service string) (id int, err error)

	//
	// Validation
	//
	HasServiceError(service string) error

	//
	// Error pages
	//
	ErrorNotFound() (code int, body []byte)
	ErrorClient(title, description string) (code int, body []byte)
	ErrorServer() (code int, body []byte)
}

// NewHTMLAdapter returns an HTMLAdapter
func NewHTMLAdapter() HTMLAdapter {
	return mock.HTMLAdapter{Code: 200, Body: []byte("not implemented yet")}
}

//
// mock.HTMLAdapter
//

// Verifies that HTMLAdapter implements mockingbird.HTMLAdapter interface
var _ HTMLAdapter = mock.HTMLAdapter{}

// NewHTMLAdapterMock returns an HTMLAdapterMock
func NewHTMLAdapterMock(code int, body string) HTMLAdapter {
	return mock.HTMLAdapter{Code: code, Body: []byte(body)}
}
