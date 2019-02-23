package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/unders/mockingbird/server/domain/mockingbird/mock"

	"github.com/unders/mockingbird/server/pkg/testdata"

	"github.com/unders/mockingbird/server/domain/mockingbird"
)

func TestAPI(t *testing.T) {
	// TODO: Add test for favicons
	t.Run("GET  /    When Content-Type=application/json    ReturnsError", getRootAsJSON)
	t.Run("GET  /    RedirectsTo  /dashboard", rootPathRedirectsToDashboard)
	t.Run("GET  /dashboard    ReturnsDashboardPage", getDashboard)
	t.Run("POST  /dashboard    ReturnsErrorRouteNotFound", postDashboard)

	t.Run("POST  /tests    StartsATestSuite    Redirects", postTests)
	t.Run("POST  /tests    WhenServerError    ReturnsError", postTestsWhenServerError)

	t.Run("GET  /tests/{id}    ReturnsTestResultPage", getTest)
	t.Run("GET  /tests/    ReturnsTestResultListPage", getTestResults)
}

func testServer(html mockingbird.HTMLAdapter) *httptest.Server {
	h := createHandler(handler{
		Favicon: func(r *http.Request) (http.Handler, bool) { return nil, false },
		HTML:    html,
		Log:     &mock.Log{},
	})

	ts := httptest.NewServer(h)
	return ts
}

func getRootAsJSON(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/",
			wantCode:       http.StatusNotImplemented,
			wantBody:       []byte(`{"error": { "code": 501, "status": "Not Implemented", "message": "JSON API not implemented by the server." } }`),
			wantRequestURL: "/",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.URL, nil)
			testdata.AssertNil(t, err)

			req.Header.Set("Accept", "application/json")
			resp, err := http.DefaultClient.Do(req)
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func rootPathRedirectsToDashboard(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("redirects to ")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL,
			wantCode:       http.StatusOK,
			wantBody:       []byte("redirects to dashboard page"),
			wantRequestURL: "/dashboard",
		},
		{
			URL:            ts.URL + "/",
			wantCode:       http.StatusOK,
			wantBody:       []byte("redirects to dashboard page"),
			wantRequestURL: "/dashboard",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Get(tc.URL)
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func getDashboard(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/dashboard",
			wantCode:       http.StatusOK,
			wantBody:       []byte("dashboard page"),
			wantRequestURL: "/dashboard",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Get(tc.URL)
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func postDashboard(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("Error: ")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/dashboard",
			wantCode:       http.StatusNotFound,
			wantBody:       []byte("Error: Not Found"),
			wantRequestURL: "/dashboard",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Post(tc.URL, "text/html", strings.NewReader("body"))
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func postTests(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("Redirects to ")})
	defer ts.Close()

	testSuiteForm := url.Values{}
	testSuiteForm.Set("test_suite", "test:all")

	testCases := []struct {
		URL            string
		contentType    string
		body           io.Reader
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/tests?test_suite=test:all",
			contentType:    "text/html",
			body:           strings.NewReader("body"),
			wantCode:       http.StatusOK,
			wantBody:       []byte("Redirects to test result page for id=test-suite-id"),
			wantRequestURL: "/tests/test-suite-id",
		},
		{
			URL:            ts.URL + "/tests",
			contentType:    "application/x-www-form-urlencoded",
			body:           strings.NewReader(testSuiteForm.Encode()),
			wantCode:       http.StatusOK,
			wantBody:       []byte("Redirects to test result page for id=test-suite-id"),
			wantRequestURL: "/tests/test-suite-id",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Post(tc.URL, tc.contentType, tc.body)
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func postTestsWhenServerError(t *testing.T) {
	serr := errors.New("")
	ts := testServer(mock.HTMLAdapter{Code: http.StatusInternalServerError, Body: []byte("Body: "), Err: serr})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/tests/?test_suite=xxx",
			wantCode:       http.StatusInternalServerError,
			wantBody:       []byte("Body: with runTest error"),
			wantRequestURL: "/tests/?test_suite=xxx",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Post(tc.URL, "text/html", strings.NewReader("body"))
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func getTest(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("body: ")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/tests/an-test-suite-id",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: test result page for id=an-test-suite-id"),
			wantRequestURL: "/tests/an-test-suite-id",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Get(tc.URL)
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}

func getTestResults(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("body: ")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/tests",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: list test result page"),
			wantRequestURL: "/tests",
		},
		{
			URL:            ts.URL + "/tests/",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: list test result page"),
			wantRequestURL: "/tests/",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Get(tc.URL)
			testdata.AssertNil(t, err)
			defer func() { testdata.AssertNil(t, resp.Body.Close()) }()

			if tc.wantCode != resp.StatusCode {
				t.Errorf("\nWant: %d\n Got: %d", tc.wantCode, resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)
			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}
		})
	}
}
