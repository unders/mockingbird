package mockingbird

// HTMLAdapter interface for working with HTML pages
type HTMLAdapter interface {
	//
	// Business  Logic
	//
	Dashboard() (code int, body []byte, err error)
	RunTest() (id string, code int, body []byte, err error)
	ShowTestResult(id string) (code int, body []byte, err error)
	ListTestResults(service string) (code int, body []byte, err error)
	RunTestForService(service string) (id string, code int, body []byte, err error)

	//
	// Error pages
	//
	ErrorNotFound() (body []byte)
	InvalidURL() (body []byte)
}
