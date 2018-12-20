package errs

import "github.com/pkg/errors"

// Forbidden error when something is forbidden
type Forbidden string

// Error error interface
func (f Forbidden) Error() string {
	return string(f)
}

// Forbidden implements forbidden interface
func (f Forbidden) Forbidden() bool {
	return true
}

// IsForbidden returns true if error is a forbidden error
//
// Usage:
//
//       if errors.IsForbidden(err) { ... }
//
func IsForbidden(err error) bool {
	type forbidden interface {
		Forbidden() bool
	}

	err = errors.Cause(err)
	if err == nil {
		return false
	}

	t, ok := err.(forbidden)
	return ok && t.Forbidden()
}
