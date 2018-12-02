package rest

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Router routes to the correct HTTP handler.
//
// Note:
//
//        If you use prefixes always set a default prefix.
//
type Router struct {
	DefaultPrefix string
	Prefixes      []string
	Namespaces    []string
}

// Route specifies a route to a specific HTTP Handler.
type Route struct {
	Method string
	Path   string
}

// New returns rest.Path and rest.Route
//
// Usage:
//
//         path, route, err := router.New(req)
//
func (r Router) New(req *http.Request) (Path, Route, error) {
	trimmedPath, err := trimPath(req)
	if err != nil {
		return Path(trimmedPath), Route{}, errors.Wrap(err, "trimPath(req) failed")
	}

	prefix, path, startIndex := adjustPrefix(trimmedPath, r.Prefixes, r.DefaultPrefix)
	startIndex = adjustStartIndexForNamespace(path, r.Namespaces, startIndex)
	rPath := routePath(path[:startIndex], path[startIndex:])

	return setPathPrefix(prefix, path), Route{Method: req.Method, Path: rPath}, nil
}

//
// PRIVATE
//

const (
	star  = "*"
	slash = "/"
	empty = ""
	colon = ":"
)

func trimPath(req *http.Request) ([]string, error) {
	path := req.URL.EscapedPath()

	var err error
	path, err = url.PathUnescape(path)
	if err != nil {
		return []string{}, errors.Wrap(err, "invalid URL format")
	}

	path = strings.TrimSuffix(path, slash)
	path = strings.TrimPrefix(path, slash)

	return strings.Split(path, slash), nil
}

// adjustPrefix adds a star instead of the prefix value from the path slice if it
// founds a match else it does nothing.
//
// Usage:
//         locale, path, startIndex := adjustPrefix([]string{"sv,", "us", "no", "en"}, "en")
//
func adjustPrefix(path, prefixes []string, defaultPrefix string) (string, []string, int) {
	if len(prefixes) == 0 {
		// When an application do not use prefixes (e.g: if it only
		// has one language) we do nothing.

		return defaultPrefix, path, 0
	}
	// else we should always have a star: * at the first position.
	size := len(path)
	pp := make([]string, size+1)
	pp[0] = star

	if size == 1 && path[0] == "" {
		return defaultPrefix, pp[:1], 1
	}

	// if prefix exist: a start replaces the prefix at first position: (len == size).
	for _, prefix := range prefixes {
		if strings.EqualFold(prefix, path[0]) {
			pp = append(pp[:1], path[1:]...)
			return prefix, pp, 1
		}
	}

	// if no prefix match: a star is inserted before the first position: (len == size + 1)
	pp = append(pp[:1], path...)
	return defaultPrefix, pp, 1
}

func adjustStartIndexForNamespace(path []string, namespaces []string, startIndex int) int {
	size := len(path)

	if startIndex >= size {
		return startIndex
	}

	if len(namespaces) == 0 {
		return startIndex
	}

	for _, namespace := range namespaces {
		if strings.EqualFold(namespace, path[startIndex]) {
			return startIndex + 1
		}
	}
	return startIndex
}

func routePath(prefix, path []string) string {
	l := len(path)

	newPrefix := append([]string{}, prefix...)

	if l == 1 && path[0] == empty {
		if len(prefix) == 0 {
			return ""
		}
		p := append(newPrefix, "")

		return slash + strings.Join(p, slash)
	}

	_path := make([]string, l)
	for i := range path {
		if i%2 == 0 {
			_path[i] = path[i]
			continue
		}

		if i == (l-1) && strings.Contains(path[i], colon) {
			operator := strings.Split(path[i], colon)[1]
			_path[i] = strings.Join([]string{star, operator}, colon)
			continue
		}

		_path[i] = star
	}

	p := append(newPrefix, _path...)
	return slash + strings.Join(p, slash)
}

func setPathPrefix(prefix string, path []string) Path {
	if prefix != empty {
		path[0] = prefix
	}

	l := len(path)
	if strings.Contains(path[l-1], colon) {
		value := strings.Split(path[l-1], colon)[0]
		path[l-1] = value
	}

	return Path(path)
}
