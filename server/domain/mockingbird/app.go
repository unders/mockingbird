package mockingbird

import (
	"fmt"
	"time"
)

type ULID string
type TestSuite string
type Status string
type State string

// The possible statuses a test suite can have
const (
	QUEUED  Status = "queued"
	RUNNING        = "running"
	DONE           = "done"
)

// The possible states a test suite can have
const (
	PENDING    State = "pending"
	SUCCESSFUL       = "successful"
	FAILED           = "failed"
)

// Dashboard shows the cumulative test suite state and possible test suites to run
type Dashboard struct {
	Stats Stats
}

// Stats shows the cumulative stats for the test suites
type Stats struct {
	LatestDoneTestSuiteID      ULID
	LatestDoneTestSuiteName    TestSuite
	LatestDoneTestSuiteState   State
	LatestDoneTestSuiteRunTime time.Duration

	LatestDoneFullTestSuiteID      ULID
	LatestDoneFullTestSuiteName    TestSuite
	LatestDoneFullTestSuiteState   State
	LatestDoneFullTestSuiteRunTime time.Duration

	TestSuiteSuccessRate float64
	TestSuiteRunCounter  int

	FullTestSuiteSuccessRate float64
	FullTestSuiteRunCounter  int

	AverageTestSuiteRunTime   time.Duration
	AverageTestSuiteQueueTime time.Duration

	SlowestTestSuiteRunTime time.Duration
	SlowestTestSuiteName    TestSuite
}

// TestResults contains a list of test results
type TestResults struct {
	NextPageToken string
	TestResults   []TestResult
}

// TestResult contains the test result for a specific test suite
type TestResult struct {
	ID ULID

	Status    Status
	State     State
	TestSuite TestSuite

	Log    []byte
	LogURL string // URL to S3 fil

	StartTime time.Time
	RunTime   time.Duration
}

func (tr TestResult) TestPath(path string) string {
	return fmt.Sprintf("%s%s", path, tr.ID)
}
func (tr TestResult) Started() string {
	return tr.StartTime.Format(time.RFC822)
}
func (tr TestResult) ShortID() string {
	return string(tr.ID)[0:10] + "..."
}
func (tr TestResult) IsPending() bool {
	return tr.State == PENDING
}
func (tr TestResult) IsSuccessful() bool {
	return tr.State == SUCCESSFUL
}
func (tr TestResult) IsFailed() bool {
	return tr.State == FAILED
}

// App defines the interface for the mockingbird application
//
// Note: it implements the business logic
//
type App interface {
	//
	// Fetches test suite results
	//
	Dashboard() (Dashboard, error)
	ListTests(pageToken string) (*TestResults, error)
	ShowTest(id ULID) (TestResult, error)
	ShowTestSuites() []TestSuite

	//
	// Executes given test suite
	//
	RunTest(s TestSuite) (ULID, error)
}
