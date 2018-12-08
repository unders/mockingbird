package mock

// Log is used for tests
type Log struct {
	ErrorLog []string
	InfoLog  []string
}

// Error records error messages
func (l *Log) Error(msg string) {
	l.ErrorLog = append(l.ErrorLog, msg)
}

// Error records info messages
func (l *Log) Info(msg string) {
	l.InfoLog = append(l.InfoLog, msg)
}
