package mock

import "github.com/unders/mockingbird/server/domain/mockingbird"

// Log is used for tests
type Log struct {
	ErrorLog []string
	InfoLog  []string
}

// Verifies that &Log implements mockingbird.Log interface
var _ mockingbird.Log = &Log{}

// Error records error messages
func (l *Log) Error(msg string) {
	l.ErrorLog = append(l.ErrorLog, msg)
}

// Error records info messages
func (l *Log) Info(msg string) {
	l.InfoLog = append(l.InfoLog, msg)
}
