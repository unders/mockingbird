package rest

import (
	"strconv"

	"github.com/pkg/errors"
)

// Path decorates a []string to make it easy to fetch values from the path.
type Path []string

// String returns string at the given position
//
// Usage:
//           value := path.String(1, "")
//
func (p Path) String(position int, defaultValue string) string {
	path := []string(p)

	if len(path) > position {
		return path[position]
	}
	return defaultValue
}

// Int returns path item at the given position as an int
func (p Path) Int(position int) (int, error) {
	value := p.String(position, "")

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return v, nil
}
