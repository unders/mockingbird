package mockingbird

import (
	"log"
)

// Log interface
type Log interface {
	Error(msg string)
	Notice(msg string)
	Info(msg string)
	Debug(msg string)
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

// Notice logs msg with Notice prefix
func (l *Logger) Notice(msg string) {
	const format = "NOTICE: %s\n"
	l.Log.Printf(format, msg)
}

// Info logs msg with INFO prefix
func (l *Logger) Info(msg string) {
	const format = "INFO: %s\n"
	l.Log.Printf(format, msg)
}

// Debug logs msg with Debug prefix
func (l *Logger) Debug(msg string) {
	const format = "Debug: %s\n"
	l.Log.Printf(format, msg)
}
