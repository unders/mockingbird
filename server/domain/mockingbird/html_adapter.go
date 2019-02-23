package mockingbird

// HTMLAdapter interface for working with HTML pages
type HTMLAdapter interface {
	//
	// Business  Logic
	//
	Dashboard() (code int, body []byte, err error)
	ListTests(pageToken string) (code int, body []byte, err error)
	ShowTest(id ULID) (code int, body []byte, err error)
	ShowTestSuites() (code int, body []byte, err error)

	RunTest(testSuite TestSuite) (id ULID, code int, body []byte, err error)

	//
	// Error pages
	//
	ErrorNotFound() (body []byte)
	InvalidURL() (body []byte)
}
