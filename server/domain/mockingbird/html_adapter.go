package mockingbird

// HTMLAdapter interface for working with HTML pages
type HTMLAdapter interface {
	//
	// Business  Logic
	//
	Dashboard() (code int, body []byte, err error)
	RunTest(environment string) (id int, err error)
	ShowTestResult(environment string, id int) (code int, body []byte, err error)
	ListTestResults(environment string) (code int, body []byte, err error)
	RunTestForService(environment, service string) (id int, err error)

	//
	// Validation
	//
	HasEnvError(environment string) error
	HasServiceError(service string) error

	//
	// Error pages
	//
	ErrorNotFound() (code int, body []byte)
	ErrorClient(title, description string) (code int, body []byte)
	ErrorServer() (code int, body []byte)
}
