package mockingbird

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
