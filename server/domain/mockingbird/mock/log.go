package mock

import "github.com/unders/mockingbird/server/domain/mockingbird"

// Log is used for tests
type Log struct {
	ErrorLog  []string
	NoticeLog []string
	InfoLog   []string
	DebugLog  []string
}

// Verifies that &mock.Log implements mockingbird.Log interface
var _ mockingbird.Log = &Log{}

// Error records error messages
func (l *Log) Error(msg string) {
	l.ErrorLog = append(l.ErrorLog, msg)
}

// Notice records notice messages
func (l *Log) Notice(msg string) {
	l.InfoLog = append(l.NoticeLog, msg)
}

// Info records info messages
func (l *Log) Info(msg string) {
	l.InfoLog = append(l.InfoLog, msg)
}

// Debug records debug messages
func (l *Log) Debug(msg string) {
	l.ErrorLog = append(l.DebugLog, msg)
}
