package app

import (
	"sync"

	"github.com/unders/mockingbird/server/pkg/uuid"

	"github.com/pkg/errors"

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
	worker     *worker
	testSuites []mockingbird.TestSuite
}

func build(ts []mockingbird.TestSuite) *Mockingbird {
	w := worker{
		work: make(chan mockingbird.ULID, 100),

		RWMutex:     &sync.RWMutex{},
		testResults: make(map[mockingbird.ULID]testResult),
	}

	// start worker queue in the background
	go func() { w.loop() }()

	return &Mockingbird{worker: &w, testSuites: ts}
}

// Verifies that *Mockingbird implements mockingbird.App interface
var _ mockingbird.App = &Mockingbird{}

//
//  Fetches test suite results
//

// Dashboard returns the Dashboard
func (m *Mockingbird) Dashboard() (mockingbird.Dashboard, error) {
	return mockingbird.Dashboard{Stats: m.worker.getStats()}, nil
}

// ListTests returns a list of test results
func (m *Mockingbird) ListTests(pageToken string) (*mockingbird.TestResults, error) {
	return m.worker.getTestResults(pageToken)
}

// ShowTest returns the test result for the given id
func (m *Mockingbird) ShowTest(id mockingbird.ULID) (mockingbird.TestResult, error) {
	return m.worker.getTestResult(id)
}

// ShowTestSuites returns the test suites
func (m *Mockingbird) ShowTestSuites() []mockingbird.TestSuite {
	return m.testSuites
}

//
// Executes a test suite
//

// RunTest executes the given test suite
func (m *Mockingbird) RunTest(s mockingbird.TestSuite) (mockingbird.ULID, error) {
	id, err := newID()
	if err != nil {
		return "", errors.Wrap(err, "newID() failed")
	}
	m.worker.add(id, s)

	return id, nil
}

func newID() (mockingbird.ULID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(err, "uuid.NewV4() failed")
	}

	return mockingbird.ULID(id), nil
}
