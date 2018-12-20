package errs

import "github.com/pkg/errors"

// NotFound error when something is not found
type NotFound string

// Error error interface
func (n NotFound) Error() string {
	return string(n)
}

// NotFound implements NotFound interface
func (n NotFound) NotFound() bool {
	return true
}

// IsNotFound returns true if error is a NotFound error
//
// Usage:
//
//       if errors.IsNotFound(err) { ... }
//
func IsNotFound(err error) bool {
	type notFound interface {
		NotFound() bool
	}

	err = errors.Cause(err)
	if err == nil {
		return false
	}

	e, ok := err.(notFound)
	return ok && e.NotFound()
}
