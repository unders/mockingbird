package mockingbird

import (
	"log"

	"github.com/unders/mockingbird/server/domain/mockingbird/mock"
)

// Log interface
type Log interface {
	Error(msg string)
	Info(msg string)
}

//
// Logger
//

// Verifies that Logger implements Log interface
var _ Log = &Logger{}

// Logger implements Log interface
type Logger struct {
	Log *log.Logger
}

// Error logs msg with ERROR prefix
func (l *Logger) Error(msg string) {
	const format = "ERROR: %s\n"
	l.Log.Printf(format, msg)
}

// Error logs msg with INFO prefix
func (l *Logger) Info(msg string) {
	const format = "INFO: %s\n"
	l.Log.Printf(format, msg)
}

//
// mock.Log
//

// Verifies that mock.Log implements mockingbird.Log interface
var _ Log = &mock.Log{}

func NewLoggerMock() Log {
	return &mock.Log{}
}
