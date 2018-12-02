package mockinbird

import "log"

// Log interface
type Log interface {
	Error(msg string)
	Info(msg string)
}

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
