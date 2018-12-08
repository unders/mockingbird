package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/unders/mockingbird/server/pkg/testdata"

	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// ******* GET  http://localhost:8080/ -> redirect-to: http://localhost:8080/v1/dashboard
// ******* GET  http://localhost:8080/v1/dashboard
//
// POST http://localhost:8080/v1/tests/
// GET  http://localhost:8080/v1/tests/{ID}
// GET  http://localhost:8080/v1/tests/?service=<service>

// POST http://localhost:8080/v1/tests/-/services/<service>

func TestAPI(t *testing.T) {
	t.Run("GET  /  RedirectsTo  /v1/dashboard", rootPathRedirectsToDashboard)
	t.Run("GET  /v1/dashboard  Returns Dashboard Page", getDashboard)
	t.Run("POST  /v1/dashboard  Returns Route Not Found  Error Page", postDashboard)
}

func testServer(t *testing.T, code int, body string) *httptest.Server {
	h := handler{
		HTML: mockingbird.NewHTMLAdapterMock(code, body),
		Log:  mockingbird.NewLoggerMock(),
	}.make()

	ts := httptest.NewServer(h)
	return ts
}

func rootPathRedirectsToDashboard(t *testing.T) {
	ts := testServer(t, 200, "redirects to ")
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL,
			wantBody:       []byte("redirects to dashboard page"),
			wantRequestURL: "/v1/dashboard",
		},
		{
			URL:            ts.URL + "/",
			wantBody:       []byte("redirects to dashboard page"),
			wantRequestURL: "/v1/dashboard",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Get(tc.URL)
			testdata.AssertNil(t, err)
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}

			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}
		})
	}
}

func getDashboard(t *testing.T) {
	ts := testServer(t, 200, "")
	defer ts.Close()

	testCases := []struct {
		URL            string
		wantBody       []byte
		wantRequestURL string
	}{
		{
			URL:            ts.URL + "/v1/dashboard",
			wantBody:       []byte("dashboard page"),
			wantRequestURL: "/v1/dashboard",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			resp, err := http.Get(tc.URL)
			testdata.AssertNil(t, err)
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			testdata.AssertNil(t, err)

			got := resp.Request.URL.RequestURI()
			if tc.wantRequestURL != got {
				t.Errorf("\nWant: %s\n Got: %s\n", tc.wantRequestURL, got)
			}

			if !reflect.DeepEqual(tc.wantBody, b) {
				t.Errorf("\nWant: %s\n Got: %s\n", string(tc.wantBody), string(b))
			}
		})
	}
}

func postDashboard(t *testing.T) {
	ts := testServer(t, 200, "Error: ")
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
			defer resp.Body.Close()

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
