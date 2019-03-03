package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/magefile/mage/sh"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/pkg/errs"
)

type testResult struct {
	mockingbird.TestResult
	index int
}

type worker struct {
	// the work buffer must be big enough :)
	work chan mockingbird.ULID

	//
	// This part should be replaced by a DB
	//
	*sync.RWMutex
	testResults map[mockingbird.ULID]testResult
	stats       mockingbird.Stats

	idIndex []mockingbird.ULID
}

func (w *worker) add(id mockingbird.ULID, s mockingbird.TestSuite) {
	tr := mockingbird.TestResult{
		ID: id,

		State:     mockingbird.PENDING,
		Status:    mockingbird.QUEUED,
		TestSuite: s,
	}

	w.Lock()
	w.idIndex = append(w.idIndex, tr.ID)
	index := len(w.idIndex) - 1
	w.testResults[id] = testResult{TestResult: tr, index: index}
	w.Unlock()

	w.work <- tr.ID
}

func (w *worker) loop() {
	for id := range w.work {
		w.workTask(id)
	}
}

func (w *worker) workTask(id mockingbird.ULID) {
	//
	// Update tr with current status and set start time
	//
	w.Lock()
	tr := w.testResults[id]
	tr.Status = mockingbird.RUNNING
	tr.StartTime = time.Now().UTC()
	w.testResults[id] = tr
	w.Unlock()

	//
	// run the test suite
	//
	out, err := sh.Output("mage", string(tr.TestSuite))
	tr.Log = out
	tr.RunTime = time.Since(tr.StartTime)
	tr.Status = mockingbird.DONE

	tr.State = mockingbird.SUCCESSFUL
	if s := sh.ExitStatus(err); s != 0 {
		tr.State = mockingbird.FAILED
		tr.Log = tr.Log + "\nerror: " + err.Error()
	}

	//
	// Save the test result
	//
	w.Lock()
	w.testResults[id] = tr
	w.Unlock()

	//
	// Update Stats
	//
	w.updateStats(tr)
}

func (w *worker) updateStats(tr testResult) {
	w.RLock()
	s := w.stats
	w.RUnlock()

	if mockingbird.FullTestSuite == tr.TestSuite {
		s.LatestDoneFullTestSuiteID = tr.ID
		s.LatestDoneFullTestSuiteName = tr.TestSuite
		s.LatestDoneFullTestSuiteState = tr.State
		s.LatestDoneFullTestSuiteRunTime = tr.RunTime

		//
		// Calculate success rate for full test suites runs
		//
		s.FullTestSuiteRunCounter = s.FullTestSuiteRunCounter + 1
		if tr.State == mockingbird.SUCCESSFUL {
			s.FullTestSuiteSuccessCounter = s.FullTestSuiteSuccessCounter + 1
		}
		s.FullTestSuiteSuccessRate = (s.FullTestSuiteSuccessCounter / s.FullTestSuiteRunCounter) * 100
	} else {
		s.LatestDoneTestSuiteID = tr.ID
		s.LatestDoneTestSuiteName = tr.TestSuite
		s.LatestDoneTestSuiteState = tr.State
		s.LatestDoneTestSuiteRunTime = tr.RunTime
	}

	//
	// Calculate success rate for all test suites
	//
	s.TestSuiteRunCounter = s.TestSuiteRunCounter + 1
	if tr.State == mockingbird.SUCCESSFUL {
		s.TestSuiteSuccessCounter = s.TestSuiteSuccessCounter + 1
	}
	s.TestSuiteSuccessRate = (s.TestSuiteSuccessCounter / s.TestSuiteRunCounter) * 100

	if tr.RunTime > s.SlowestTestSuiteRunTime {
		s.SlowestTestSuiteName = tr.TestSuite
		s.SlowestTestSuiteRunTime = tr.RunTime
	}

	w.Lock()
	w.stats = s
	w.Unlock()
}

func (w *worker) getStats() mockingbird.Stats {
	return w.stats
}

func (w *worker) getTestResults(pageToken string) (*mockingbird.TestResults, error) {
	startIndex := len(w.idIndex) - 1
	if pageToken != "" {
		if ts, ok := w.testResults[mockingbird.ULID(pageToken)]; ok {
			startIndex = ts.index
		}
	}

	nextPageToken := ""
	stopIndex := startIndex - 10
	nextPageTokenIndex := startIndex - 11
	if stopIndex < 1 {
		stopIndex = 0
		nextPageTokenIndex = -1
	}

	var trs []mockingbird.TestResult

	if nextPageTokenIndex > -1 {
		id := w.idIndex[nextPageTokenIndex]
		if tr, ok := w.testResults[id]; ok {
			nextPageToken = string(tr.ID)
		}
	}

	for i := startIndex; i >= stopIndex; i-- {
		id := w.idIndex[i]
		tr, ok := w.testResults[id]
		if !ok {
			return nil, errors.Errorf("w.testResults[%s] failed", string(id))
		}
		trs = append(trs, tr.TestResult)
	}

	return &mockingbird.TestResults{NextPageToken: nextPageToken, TestResults: trs}, nil
}

func (w *worker) getTestResult(id mockingbird.ULID) (mockingbird.TestResult, error) {
	w.RLock()
	defer w.RUnlock()

	if ts, ok := w.testResults[id]; ok {
		return ts.TestResult, nil
	}

	msg := fmt.Sprintf("test result %s not found", id)
	return mockingbird.TestResult{}, errs.NotFound(msg)
}
