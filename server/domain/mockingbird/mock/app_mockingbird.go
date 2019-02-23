package mock

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/pkg/errs"
)

// ULIDs
const (
	ULID1        mockingbird.ULID = "01BX5ZZKBKACTAV9WEVGEMMVRY" // test:registration
	ULID2        mockingbird.ULID = "01BX5ZZKBKACTAV9WEVGEMMVRZ" // test:navigation
	ULID3        mockingbird.ULID = "01BX5ZZKBKACTAV9WEVGEMMVS0" // test:all
	NotFoundULID mockingbird.ULID = "01BX5ZZKBKZZZZZZZZZZZZZZZZ" // not connected to any test suite
)

// test suites
const (
	TestAll          mockingbird.TestSuite = "all:test"
	TestNavigation                         = "navigation:test"
	TestRegistration                       = "registration:test"
)

// test suite list
var (
	testSuites = []mockingbird.TestSuite{
		TestAll,
		TestNavigation,
		TestRegistration,
	}
)

// AppMockingbird implements the mockingbird.App interface
//
// Note:
//
//        mock.AppMockingbird simulate the business logic of the real application
//
//
type AppMockingbird struct {
	Now time.Time
}

// Verifies that *mock.AppMockingbird implements mockingbird.App interface
var _ mockingbird.App = &AppMockingbird{}

//
//  Fetches test suite results
//

// Dashboard returns the Dashboard
func (m *AppMockingbird) Dashboard() (mockingbird.Dashboard, error) {
	d := mockingbird.Dashboard{
		Stats: mockingbird.Stats{
			LatestDoneFullTestSuiteID:   ULID3,
			LatestDoneFullTestSuiteName: TestAll,
			LatestDoneTestSuiteState:    mockingbird.SUCCESSFUL,
			LatestDoneTestSuiteRunTime:  66 * time.Second,

			LatestDoneTestSuiteID:          ULID1,
			LatestDoneTestSuiteName:        TestRegistration,
			LatestDoneFullTestSuiteState:   mockingbird.FAILED,
			LatestDoneFullTestSuiteRunTime: 150 * time.Second,

			TestSuiteSuccessRate: 90,
			TestSuiteRunCounter:  17931,

			FullTestSuiteSuccessRate: 99,
			FullTestSuiteRunCounter:  9840,

			AverageTestSuiteRunTime:   55 * time.Second,
			AverageTestSuiteQueueTime: 106 * time.Second,

			SlowestTestSuiteRunTime: 255 * time.Second,
			SlowestTestSuiteName:    TestAll,
		},
	}
	return d, nil
}

// ListTests returns a list of test results
func (m *AppMockingbird) ListTests(pageToken string) (*mockingbird.TestResults, error) {
	if pageToken == "error" {
		return nil, errors.New("server error")
	}

	return m.listTests(), nil
}

// ShowTest returns the test result for the given id
func (m *AppMockingbird) ShowTest(id mockingbird.ULID) (mockingbird.TestResult, error) {
	for _, test := range m.listTests().TestResults {
		if test.ID == id {
			return test, nil
		}
	}

	msg := fmt.Sprintf("test run %s not found", id)
	return mockingbird.TestResult{}, errs.NotFound(msg)
}

// ShowTestSuite returns the test suites
func (m *AppMockingbird) ShowTestSuites() []mockingbird.TestSuite {
	return testSuites
}

//
// Executes a test suite
//

// RunTest executes the given test suite
func (m *AppMockingbird) RunTest(s mockingbird.TestSuite) (mockingbird.ULID, error) {
	switch s {
	case TestRegistration:
		return ULID1, nil
	case TestNavigation:
		return ULID2, nil
	case TestAll:
		return ULID3, nil
	default:
		msg := fmt.Sprintf("Invalid test suite name: `%s`", s)
		return "", errs.NotFound(msg)
	}
}

//
// PRIVATE
//

func (m *AppMockingbird) listTests() *mockingbird.TestResults {
	return &mockingbird.TestResults{
		NexPageToken: string(ULID1),
		TestResults: []mockingbird.TestResult{
			{
				ID: ULID1,

				Status:    mockingbird.DONE,
				State:     mockingbird.SUCCESSFUL,
				TestSuite: TestRegistration,

				Log:    testAllLogOutput,
				LogURL: "...", // An URL to an S3 file

				StartTime: m.Now,
				RunTime:   100 * time.Second,
			},
			{
				ID: ULID2,

				Status:    mockingbird.RUNNING,
				State:     mockingbird.PENDING,
				TestSuite: TestNavigation,

				Log:    testAllLogOutput,
				LogURL: "...", // An URL to an S3 file

				StartTime: m.Now,
				RunTime:   100 * time.Second,
			},
			{
				ID: ULID3,

				Status:    mockingbird.QUEUED,
				State:     mockingbird.PENDING,
				TestSuite: TestAll,

				Log:    testAllLogOutput,
				LogURL: "...", // An URL to an S3 file

				StartTime: m.Now,
				RunTime:   100 * time.Second,
			},
		},
	}
}

var testAllLogOutput = []byte(`
=== RUN   TestSearch
--- PASS: TestSearch (0.00s)
    search_test.go:9: doing the test...
PASS
ok  	github.com/unders/mockingbird/test/service/google	0.006s
`)
