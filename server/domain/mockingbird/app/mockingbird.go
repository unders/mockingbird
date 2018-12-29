package app

import (
	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// Mockingbird implements the mockingbird.App interface
//
// Note:
//
//        app.Mockingbird implements the business logic of the application
//
//
type Mockingbird struct {
}

// Verifies that *Mockingbird implements mockingbird.App interface
var _ mockingbird.App = &Mockingbird{}

//
//  Fetches test suite results
//

// Dashboard returns the Dashboard
func (m *Mockingbird) Dashboard() (mockingbird.Dashboard, error) {
	return mockingbird.Dashboard{}, nil
}

// ListTests returns a list of test results
func (m *Mockingbird) ListTests(pageToken string) (*mockingbird.TestResults, error) {
	return &mockingbird.TestResults{}, nil
}

// ShowTest returns the test result for the given id
func (m *Mockingbird) ShowTest(id mockingbird.ULID) (mockingbird.TestResult, error) {
	return mockingbird.TestResult{}, nil
}

//
// Executes a test suite
//

// RunTest executes the given test suite
func (m *Mockingbird) RunTest(s mockingbird.TestSuite) (mockingbird.ULID, error) {
	return "", nil
}
