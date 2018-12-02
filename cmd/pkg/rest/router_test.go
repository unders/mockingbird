package rest_test

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/unders/mockingbird/cmd/pkg/testdata"

	"github.com/unders/mockingbird/cmd/pkg/rest"
)

func TestMatcher_NewRoute_WithNamespaceAndPrefix(t *testing.T) {
	router := rest.Router{
		Namespaces:    []string{"v2", "admin", "v3"},
		Prefixes:      []string{"sv", "en", "no"},
		DefaultPrefix: "sv",
	}
	tests := []struct {
		req   *http.Request
		route rest.Route
		path  []string
	}{
		{
			req:   newRequest(t, http.MethodPost, ""),
			route: rest.Route{http.MethodPost, "/*"},
			path:  []string{"sv"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/"),
			route: rest.Route{http.MethodGet, "/*"},
			path:  []string{"sv"},
		},
		{
			req:   newRequest(t, http.MethodPost, "admin"),
			route: rest.Route{http.MethodPost, "/*/admin"},
			path:  []string{"sv", "admin"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/en/admin/shelves/"),
			route: rest.Route{http.MethodPost, "/*/admin/shelves"},
			path:  []string{"en", "admin", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodPost, "admin/shelves/"),
			route: rest.Route{http.MethodPost, "/*/admin/shelves"},
			path:  []string{"sv", "admin", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodPatch, "en/admin/shelves/100"),
			route: rest.Route{http.MethodPatch, "/*/admin/shelves/*"},
			path:  []string{"en", "admin", "shelves", "100"},
		},
		{
			req:   newRequest(t, http.MethodPatch, "admin/shelves/100"),
			route: rest.Route{http.MethodPatch, "/*/admin/shelves/*"},
			path:  []string{"sv", "admin", "shelves", "100"},
		},
		{
			req:   newRequest(t, http.MethodPatch, "en/admin/shelves/100/books"),
			route: rest.Route{http.MethodPatch, "/*/admin/shelves/*/books"},
			path:  []string{"en", "admin", "shelves", "100", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/en/admin/shelves/100/books:clear"),
			route: rest.Route{http.MethodPost, "/*/admin/shelves/*/books:clear"},
			path:  []string{"en", "admin", "shelves", "100", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "en/admin/shelves/-/books/100:cancel"),
			route: rest.Route{http.MethodPost, "/*/admin/shelves/*/books/*:cancel"},
			path:  []string{"en", "admin", "shelves", "-", "books", "100"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/users/john%20smith/events/123"),
			route: rest.Route{http.MethodGet, "/*/v3/users/*/events/*"},
			path:  []string{"sv", "v3", "users", "john smith", "events", "123"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/en/v3/users/john%20smith/events/123"),
			route: rest.Route{http.MethodGet, "/*/v3/users/*/events/*"},
			path:  []string{"en", "v3", "users", "john smith", "events", "123"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/events:clear"),
			route: rest.Route{http.MethodGet, "/*/v3/events:clear"},
			path:  []string{"sv", "v3", "events"},
		},

		{
			req:   newRequest(t, http.MethodPost, "/v2/events/100:cancel"),
			route: rest.Route{http.MethodPost, "/*/v2/events/*:cancel"},
			path:  []string{"sv", "v2", "events", "100"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/en/v2/events/100:cancel"),
			route: rest.Route{http.MethodPost, "/*/v2/events/*:cancel"},
			path:  []string{"en", "v2", "events", "100"},
		},
		{
			req:   newRequest(t, http.MethodPost, "xxx/v2/events/100:cancel"),
			route: rest.Route{http.MethodPost, "/*/xxx/*/events/*:cancel"},
			path:  []string{"sv", "xxx", "v2", "events", "100"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/en/v2/shelves/-/books"),
			route: rest.Route{http.MethodGet, "/*/v2/shelves/*/books"},
			path:  []string{"en", "v2", "shelves", "-", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "v3/shelves/1/books/67:cancel"),
			route: rest.Route{http.MethodPost, "/*/v3/shelves/*/books/*:cancel"},
			path:  []string{"sv", "v3", "shelves", "1", "books", "67"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/no/v2/shelves/1/books:clear"),
			route: rest.Route{http.MethodPost, "/*/v2/shelves/*/books:clear"},
			path:  []string{"no", "v2", "shelves", "1", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/no/shelves/1/books:clear"),
			route: rest.Route{http.MethodPost, "/*/shelves/*/books:clear"},
			path:  []string{"no", "shelves", "1", "books"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			path, route, err := router.New(test.req)
			testdata.AssertNil(t, err)

			if !reflect.DeepEqual(test.route, route) {
				t.Errorf("\nWant: %+v\n Got: %+v\n", test.route, route)
			}
			testPath(t, test.path, path)
		})
	}
}

func TestMatcher_NewRoute_WithNamespace(t *testing.T) {
	router := rest.Router{Namespaces: []string{"v2", "v3"}}
	tests := []struct {
		req   *http.Request
		route rest.Route
		path  []string
	}{
		{
			req:   newRequest(t, http.MethodPost, ""),
			route: rest.Route{http.MethodPost, ""},
			path:  []string{""},
		},
		{
			req:   newRequest(t, http.MethodGet, "/"),
			route: rest.Route{http.MethodGet, ""},
			path:  []string{""},
		},
		{
			req:   newRequest(t, http.MethodPost, "v3"),
			route: rest.Route{http.MethodPost, "/v3"},
			path:  []string{"v3"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3"),
			route: rest.Route{http.MethodGet, "/v3"},
			path:  []string{"v3"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/shelves"),
			route: rest.Route{http.MethodGet, "/v3/shelves"},
			path:  []string{"v3", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/shelves/100"),
			route: rest.Route{http.MethodGet, "/v3/shelves/*"},
			path:  []string{"v3", "shelves", "100"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/shelves/100/books"),
			route: rest.Route{http.MethodGet, "/v3/shelves/*/books"},
			path:  []string{"v3", "shelves", "100", "books"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/shelves/-/books/9000"),
			route: rest.Route{http.MethodGet, "/v3/shelves/*/books/*"},
			path:  []string{"v3", "shelves", "-", "books", "9000"},
		},

		{
			req:   newRequest(t, http.MethodGet, "/v3/users/john%20smith/events/123"),
			route: rest.Route{http.MethodGet, "/v3/users/*/events/*"},
			path:  []string{"v3", "users", "john smith", "events", "123"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v3/events:clear"),
			route: rest.Route{http.MethodGet, "/v3/events:clear"},
			path:  []string{"v3", "events"},
		},

		{
			req:   newRequest(t, http.MethodPost, "/v2/events/100:cancel"),
			route: rest.Route{http.MethodPost, "/v2/events/*:cancel"},
			path:  []string{"v2", "events", "100"},
		},
		{
			req:   newRequest(t, http.MethodPost, "xxx/v2/events/100:cancel"),
			route: rest.Route{http.MethodPost, "/xxx/*/events/*:cancel"},
			path:  []string{"xxx", "v2", "events", "100"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/v2/shelves/-/books"),
			route: rest.Route{http.MethodGet, "/v2/shelves/*/books"},
			path:  []string{"v2", "shelves", "-", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "v3/shelves/1/books/67:cancel"),
			route: rest.Route{http.MethodPost, "/v3/shelves/*/books/*:cancel"},
			path:  []string{"v3", "shelves", "1", "books", "67"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/v2/shelves/1/books:clear"),
			route: rest.Route{http.MethodPost, "/v2/shelves/*/books:clear"},
			path:  []string{"v2", "shelves", "1", "books"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			path, route, err := router.New(test.req)
			testdata.AssertNil(t, err)

			if !reflect.DeepEqual(test.route, route) {
				t.Errorf("\nWant: %+v\n Got: %+v\n", test.route, route)
			}
			testPath(t, test.path, path)
		})
	}
}
func TestMatcher_NewRoute_WithPrefix(t *testing.T) {
	router := rest.Router{Prefixes: []string{"sv", "en", "no"}, DefaultPrefix: "sv"}
	tests := []struct {
		req   *http.Request
		route rest.Route
		path  []string
	}{
		{
			// when no prefix is present in path, the default prefix sv is used
			req:   newRequest(t, http.MethodPost, ""),
			route: rest.Route{http.MethodPost, "/*"},
			path:  []string{"sv"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/"),
			route: rest.Route{http.MethodGet, "/*"},
			path:  []string{"sv"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/shelves/"),
			route: rest.Route{http.MethodGet, "/*/shelves"},
			path:  []string{"sv", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/shelves"),
			route: rest.Route{http.MethodGet, "/*/shelves"},
			path:  []string{"sv", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/sv"),
			route: rest.Route{http.MethodPost, "/*"},
			path:  []string{"sv"},
		},
		{
			req:   newRequest(t, http.MethodGet, "sv"),
			route: rest.Route{http.MethodGet, "/*"},
			path:  []string{"sv"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/sv/shelves"),
			route: rest.Route{http.MethodGet, "/*/shelves"},
			path:  []string{"sv", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/en"),
			route: rest.Route{http.MethodPost, "/*"},
			path:  []string{"en"},
		},
		{
			req:   newRequest(t, http.MethodGet, "en"),
			route: rest.Route{http.MethodGet, "/*"},
			path:  []string{"en"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/en/shelves"),
			route: rest.Route{http.MethodGet, "/*/shelves"},
			path:  []string{"en", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/xxx/shelves"),
			route: rest.Route{http.MethodGet, "/*/xxx/*"},
			path:  []string{"sv", "xxx", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/shelves"),
			route: rest.Route{http.MethodPost, "/*/shelves"},
			path:  []string{"sv", "shelves"},
		},
		{
			req:   newRequest(t, http.MethodPut, "no/shelves/900"),
			route: rest.Route{http.MethodPut, "/*/shelves/*"},
			path:  []string{"no", "shelves", "900"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/shelves/900/books"),
			route: rest.Route{http.MethodGet, "/*/shelves/*/books"},
			path:  []string{"sv", "shelves", "900", "books"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/no/shelves/-/books"),
			route: rest.Route{http.MethodGet, "/*/shelves/*/books"},
			path:  []string{"no", "shelves", "-", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "en/shelves/1/books/67:cancel"),
			route: rest.Route{http.MethodPost, "/*/shelves/*/books/*:cancel"},
			path:  []string{"en", "shelves", "1", "books", "67"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/no/shelves/1/books:clear"),
			route: rest.Route{http.MethodPost, "/*/shelves/*/books:clear"},
			path:  []string{"no", "shelves", "1", "books"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			path, route, err := router.New(test.req)
			testdata.AssertNil(t, err)

			if !reflect.DeepEqual(test.route, route) {
				t.Errorf("\nWant: %+v\n Got: %+v\n", test.route, route)
			}
			testPath(t, test.path, path)
		})
	}
}

func TestMatcher_NewRoute_WhenNoNamespaceOrPrefix(t *testing.T) {
	router := rest.Router{}
	tests := []struct {
		req   *http.Request
		route rest.Route
		path  []string
	}{
		{
			req:   newRequest(t, http.MethodPost, ""),
			route: rest.Route{http.MethodPost, ""},
			path:  []string{""},
		},
		{
			req:   newRequest(t, http.MethodGet, "/"),
			route: rest.Route{http.MethodGet, ""},
			path:  []string{""},
		},
		{
			req:   newRequest(t, http.MethodPost, "/shelves"),
			route: rest.Route{http.MethodPost, "/shelves"},
			path:  []string{"shelves"},
		},
		{
			req:   newRequest(t, http.MethodPut, "/shelves/900"),
			route: rest.Route{http.MethodPut, "/shelves/*"},
			path:  []string{"shelves", "900"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/shelves/900/books"),
			route: rest.Route{http.MethodGet, "/shelves/*/books"},
			path:  []string{"shelves", "900", "books"},
		},
		{
			req:   newRequest(t, http.MethodGet, "/shelves/-/books"),
			route: rest.Route{http.MethodGet, "/shelves/*/books"},
			path:  []string{"shelves", "-", "books"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/shelves/1/books/67:cancel"),
			route: rest.Route{http.MethodPost, "/shelves/*/books/*:cancel"},
			path:  []string{"shelves", "1", "books", "67"},
		},
		{
			req:   newRequest(t, http.MethodPost, "/shelves/1/books:clear"),
			route: rest.Route{http.MethodPost, "/shelves/*/books:clear"},
			path:  []string{"shelves", "1", "books"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			path, route, err := router.New(test.req)
			testdata.AssertNil(t, err)

			if !reflect.DeepEqual(test.route, route) {
				t.Errorf("\nWant: %+v\n Got: %+v\n", test.route, route)
			}
			testPath(t, test.path, path)
		})
	}
}

func testPath(t *testing.T, wantPath []string, gotPath rest.Path) {
	t.Helper()

	if len(wantPath) != len(gotPath) {
		t.Fatalf("\nWant: %d\n Got: %d\nPath: %+v\n", len(wantPath), len(gotPath), gotPath)
	}

	for i, want := range wantPath {
		if !strings.EqualFold(want, gotPath[i]) {
			t.Errorf("\nWant: %s\n Got: %s\npath: %+v\n", want, gotPath[i], gotPath)
		}
	}
}

func newRequest(t *testing.T, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	testdata.AssertNil(t, err)
	return req
}
