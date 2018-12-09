package mockingbird

import "github.com/unders/mockingbird/server/domain/mockingbird/mock"

// HTMLAdapter interface for working with HTML pages
type HTMLAdapter interface {
	//
	// Business  Logic
	//
	Dashboard() (code int, body []byte, err error)
	RunTest() (id string, code int, body []byte, err error)
	ShowTestResult(id string) (code int, body []byte, err error)
	ListTestResults() (code int, body []byte, err error)
	RunTestForService(service string) (id string, code int, body []byte, err error)

	//
	// Validation
	//
	HasIdError(id string) (code int, body []byte, err error)
	HasServiceError(service string) (code int, body []byte, err error)

	//
	// Error pages
	//
	ErrorNotFound() (body []byte)
	InvalidURL() (body []byte)
}

// NewHTMLAdapter returns an HTMLAdapter
func NewHTMLAdapter() HTMLAdapter {
	return mock.HTMLAdapter{Code: 200, Body: []byte("not implemented yet")}
}

//
// mock.HTMLAdapter
//

// Verifies that mock.HTMLAdapter implements mockingbird.HTMLAdapter interface
var _ HTMLAdapter = mock.HTMLAdapter{}

// NewHTMLAdapterMock returns an HTMLAdapterMock
func NewHTMLAdapterMock(code int, body string, err, idErr, serviceErr error) HTMLAdapter {
	return mock.HTMLAdapter{
		Code:       code,
		Body:       []byte(body),
		Err:        err,
		IDErr:      idErr,
		ServiceErr: serviceErr,
	}
}
