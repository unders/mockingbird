package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	t.Run("GET  /  When Content-Type=application/json ReturnsError", getRootAsJSON)
	t.Run("GET  /  RedirectsTo  /v1/dashboard", rootPathRedirectsToDashboard)
	t.Run("GET  /v1/dashboard  ReturnsDashboardPage", getDashboard)
	t.Run("POST  /v1/dashboard  ReturnsErrorRouteNotFound", postDashboard)

	t.Run("POST  /v1/tests  StartsATestSuite  Redirects", postTests)
	t.Run("POST  /v1/tests  WhenServerError  ReturnsError", postTestsWhenServerError)

	t.Run("GET  /v1/tests/{id} ReturnsTestResultPage", getTest)
	t.Run("GET  /v1/tests/ ReturnsTestResultListPage", getTestResults)

	t.Run("POST  /v1/tests/-/services/{service-name} StartsATestSuite  Redirects", postServiceTests)
	t.Run("POST  /v1/tests/-/services/{service-name} WhenInvalidServiceName  ReturnsError", postServiceTestsWhenInvalidService)
	t.Run("POST  /v1/tests/-/services/{service-name} WhenServerError  ReturnsError", postServiceTestsWhenServerError)
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
			wantRequestURL: "/v1/dashboard",
		},
		{
			URL:            ts.URL + "/",
			wantCode:       http.StatusOK,
			wantBody:       []byte("redirects to dashboard page"),
			wantRequestURL: "/v1/dashboard",
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
			URL:            ts.URL + "/v1/dashboard",
			wantCode:       http.StatusOK,
			wantBody:       []byte("dashboard page"),
			wantRequestURL: "/v1/dashboard",
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
			URL:            ts.URL + "/v1/dashboard",
			wantCode:       http.StatusBadRequest,
			wantBody:       []byte("Error: Not Found"),
			wantRequestURL: "/v1/dashboard",
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

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/v1/tests",
			wantCode:       http.StatusOK,
			wantBody:       []byte("Redirects to test result page for id=test-suite-id"),
			wantRequestURL: "/v1/tests/test-suite-id",
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
			URL:            ts.URL + "/v1/tests/",
			wantCode:       http.StatusInternalServerError,
			wantBody:       []byte("Body: with runTest error"),
			wantRequestURL: "/v1/tests/",
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
			URL:            ts.URL + "/v1/tests/an-test-suite-id",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: test result page for id=an-test-suite-id"),
			wantRequestURL: "/v1/tests/an-test-suite-id",
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
			URL:            ts.URL + "/v1/tests",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: list test result page"),
			wantRequestURL: "/v1/tests",
		},
		{
			URL:            ts.URL + "/v1/tests/",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: list test result page"),
			wantRequestURL: "/v1/tests/",
		},
		{
			URL:            ts.URL + "/v1/tests/?service=service-name1",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: list test result page for service=service-name1"),
			wantRequestURL: "/v1/tests/?service=service-name1",
		},
		{
			URL:            ts.URL + "/v1/tests?service=service-name2",
			wantCode:       http.StatusOK,
			wantBody:       []byte("body: list test result page for service=service-name2"),
			wantRequestURL: "/v1/tests?service=service-name2",
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

func postServiceTests(t *testing.T) {
	ts := testServer(mock.HTMLAdapter{Code: http.StatusOK, Body: []byte("Redirects to ")})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/v1/tests/-/services/service-x",
			wantCode:       http.StatusOK,
			wantBody:       []byte("Redirects to test result page for id=test-suite-x"),
			wantRequestURL: "/v1/tests/test-suite-x",
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

func postServiceTestsWhenInvalidService(t *testing.T) {
	serr := errors.New("service name is invalid")
	ts := testServer(mock.HTMLAdapter{Code: http.StatusBadRequest, Body: []byte("Body: "), ServiceErr: serr})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/v1/tests/-/services/service-x",
			wantCode:       http.StatusBadRequest,
			wantBody:       []byte("Body: HasServiceError for service-x"),
			wantRequestURL: "/v1/tests/-/services/service-x",
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

func postServiceTestsWhenServerError(t *testing.T) {
	serr := errors.New("server error")
	ts := testServer(mock.HTMLAdapter{Code: http.StatusInternalServerError, Body: []byte("Body: "), Err: serr})
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantCode       int
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/v1/tests/-/services/service-x",
			wantCode:       http.StatusInternalServerError,
			wantBody:       []byte("Body: with runTestForService error service=service-x"),
			wantRequestURL: "/v1/tests/-/services/service-x",
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
