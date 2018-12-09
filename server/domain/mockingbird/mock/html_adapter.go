package mock

// HTMLAdapter is used for tests
type HTMLAdapter struct {
	Code       int
	Body       []byte
	Err        error
	IDErr      error
	ServiceErr error
}

//
// Business  Logic
//

// Dashboard returns the dashboard page
func (a HTMLAdapter) Dashboard() (code int, body []byte, err error) {
	b := append(a.Body, "dashboard page"...)
	if a.Err != nil {
		b = append(b, " with error"...)
	}
	return a.Code, b, a.Err
}

// RunTest starts a test suite
func (a HTMLAdapter) RunTest() (id string, code int, cody []byte, err error) {
	b := a.Body
	if a.Err != nil {
		b = append(b, " with runTest error"...)
	}
	return "test-suite-id", a.Code, b, a.Err
}

// ShowTestResult returns the ShowTestResult page
func (a HTMLAdapter) ShowTestResult(id string) (code int, body []byte, err error) {
	b := append(a.Body, "test result page for id="...)
	b = append(b, id...)
	if a.Err != nil {
		b = append(b, " with error"...)
	}
	return a.Code, b, a.Err
}

// ListTestResults returns the ListTestResults page
func (a HTMLAdapter) ListTestResults() (code int, body []byte, err error) {
	b := append(a.Body, "list test result page"...)
	if a.Err != nil {
		b = append(b, " with error"...)
	}
	return a.Code, b, a.Err
}

// RunTestForService starts a test suite for a service
func (a HTMLAdapter) RunTestForService(service string) (id string, code int, body []byte, err error) {
	b := a.Body
	if a.Err != nil {
		b = append(b, " with runTestForService error"...)
		b = append(b, service...)
	}
	return "test-suite-for-service-id", a.Code, b, a.Err
}

//
// Validations
//
func (a HTMLAdapter) HasIdError(id string) (code int, body []byte, err error) {
	if a.IDErr != nil {
		b := append(a.Body, "HasIdError for "...)
		b = append(b, id...)

		return a.Code, b, a.IDErr
	}
	return a.Code, a.Body, nil
}

func (a HTMLAdapter) HasServiceError(service string) (code int, body []byte, err error) {
	if a.ServiceErr != nil {
		b := append(a.Body, "HasServiceError for "...)
		b = append(b, service...)

		return a.Code, b, a.ServiceErr
	}
	return a.Code, a.Body, nil
}

//
// Error pages
//

// ErrorNotFound returns error not found page
func (a HTMLAdapter) ErrorNotFound() (body []byte) {
	b := append(a.Body, "Not Found"...)
	return b
}

// InvalidURL returns error an Invalid URL page
func (a HTMLAdapter) InvalidURL() (body []byte) {
	b := append(a.Body, "Invalid URL page"...)
	return b
}
