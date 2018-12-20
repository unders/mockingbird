package mockingbird

// App defines the interface for the mockingbird application
//
// Note: it implements the business logic
//
type App interface {
	Dashboard() error
	RunTest() (id string, err error)
	ShowTestResult(id string) error
	ListTestResults(service string) error
	RunTestForService(service string) (id string, err error)
}
