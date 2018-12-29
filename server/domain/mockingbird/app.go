package mockingbird

import "time"

type ULID string
type TestSuite string

type status string
type state string

// The possible statuses a test suite can have
const (
	QUEUED  status = "queued"
	RUNNING        = "running"
	DONE           = "done"
)

// The possible states a test suite can have
const (
	PENDING    state = "pending"
	SUCCESSFUL       = "successful"
	FAILED           = "failed"
)

// Dashboard shows the cumulative test suite state and possible test suites to run
type Dashboard struct {
	Stats      Stats
	TestSuites []TestSuite
}

// Stats shows the cumulative stats for the test suites
type Stats struct {
	LatestDoneTestSuiteState   state
	LatestDoneTestSuiteRunTime time.Time

	LatestDoneFullTestSuiteState   state
	LatestDoneFullTestSuiteRunTime time.Time

	TestSuiteSuccessRate float64
	TestSuiteRunCounter  int

	FullTestSuiteSuccessRate float64
	FullTesSuiteRunCounter   int

	AverageTestSuiteRunTime   time.Duration
	AverageTestSuiteQueueTime time.Duration
	SlowestTestSuiteRunTime   time.Duration

	SlowestTestSuiteName TestSuite
}

// TestResults contains a list of test results
type TestResults struct {
	NexPageToken string
	TestResults  []TestResult
}

// TestResult contains the test result for a specific test suite
type TestResult struct {
	ID ULID

	Status    status
	State     state
	TestSuite TestSuite

	Log    []byte
	LogURL string // URL to S3 fil

	StartTime time.Time
	RunTime   time.Time
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

	//
	// Executes given test suite
	//
	RunTest(s TestSuite) (ULID, error)
}
