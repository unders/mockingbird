package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/unders/mockingbird/server/domain/mockingbird"

	"github.com/pkg/errors"
	"github.com/unders/mockingbird/server/pkg/rest"
)

type handler struct {
	HTML mockingbird.HTMLAdapter
	Log  mockingbird.Log
}

func (h handler) make() http.Handler {
	router := rest.Router{Namespaces: []string{"v1"}}

	f := func(w http.ResponseWriter, req *http.Request) {
		if h.isJSON(req) {
			h.jsonNotImplemented(w, req)
			return
		}
		//
		// All other Content-Type's are treated as HTML.
		//

		path, route, err := router.New(req)
		if err != nil {
			h.write(w, req, http.StatusBadRequest, h.HTML.InvalidURL(), err)
			return
		}

		switch route {
		case rest.Route{Method: http.MethodGet, Path: ""}:
			url := "/v1/dashboard"
			http.Redirect(w, req, url, http.StatusSeeOther)
		case rest.Route{Method: http.MethodGet, Path: "/v1/dashboard"}:
			h.showDashboard(w, req)
		case rest.Route{Method: http.MethodPost, Path: "/v1/tests"}:
			h.runTest(w, req)
		case rest.Route{Method: http.MethodGet, Path: "/v1/tests/*"}:
			h.showTestResult(w, req, path)
		case rest.Route{Method: http.MethodGet, Path: "/v1/tests"}:
			h.listTestResults(w, req)
		case rest.Route{Method: http.MethodPost, Path: "/v1/tests/-/services/*"}:
			h.runTestForService(w, req, path)
		default:
			err := errors.New("route not found")
			h.write(w, req, http.StatusBadRequest, h.HTML.ErrorNotFound(), err)
		}
	}

	return http.HandlerFunc(f)
}

//
// HTML Handlers
//
func (h *handler) showDashboard(w http.ResponseWriter, req *http.Request) {
	code, b, err := h.HTML.Dashboard()
	h.write(w, req, code, b, err)
}

func (h *handler) runTest(w http.ResponseWriter, req *http.Request) {
	id, code, body, err := h.HTML.RunTest()
	if err != nil {
		h.write(w, req, code, body, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/v1/tests/%s", id), http.StatusSeeOther)
}

func (h *handler) showTestResult(w http.ResponseWriter, req *http.Request, path rest.Path) {
	code, b, err := h.HTML.ShowTestResult(path.String(2, ""))
	h.write(w, req, code, b, err)
}

func (h *handler) listTestResults(w http.ResponseWriter, req *http.Request) {
	code, b, err := h.HTML.ListTestResults()
	h.write(w, req, code, b, err)
}

func (h *handler) runTestForService(w http.ResponseWriter, req *http.Request, path rest.Path) {
	id, code, body, err := h.HTML.RunTestForService(path.String(4, ""))
	if err != nil {
		h.write(w, req, code, body, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/v1/tests/%s", id), http.StatusSeeOther)
}

//
// HTML Writer
//

func (h *handler) write(w http.ResponseWriter, req *http.Request, code int, buf []byte, err error) {
	body := bytes.NewReader(buf)

	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	w.WriteHeader(code)
	if _, err := io.Copy(w, body); err != nil {
		errors.Wrapf(err, "io.Copy(w, body) failed")
		h.logResponseFailure(req, code, err)
		return
	}

	h.logRequest(req, code, err)
}

//
// JSON
//

func (h *handler) isJSON(req *http.Request) bool {
	ct := req.Header.Get("Content-Type")
	return ct == "application/json"
}

func (h *handler) jsonNotImplemented(w http.ResponseWriter, req *http.Request) {
	const code = http.StatusNotImplemented
	var status = http.StatusText(code)
	const msg = "JSON API not implemented by the server."

	const format = `{"error": { "code": %d, "status": "%s", "message": "%s" } }`
	b := []byte(fmt.Sprintf(format, code, status, msg))
	err := errors.Errorf(format, code, status, msg)
	h.writeJSON(w, req, code, b, err)
}

func (h *handler) writeJSON(w http.ResponseWriter, req *http.Request, code int, buf []byte, err error) {
	body := bytes.NewReader(buf)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	w.WriteHeader(code)
	if _, err := io.Copy(w, body); err != nil {
		errors.Wrap(err, "io.Copy(w, body) failed")
		h.logResponseFailure(req, code, err)
		return
	}

	h.logRequest(req, code, err)
}

//
// General request logging
//
func (h *handler) logResponseFailure(req *http.Request, code int, err error) {
	msg := http.StatusText(code)
	const format = "%s %s    response error=%s    [%d %s]"
	h.Log.Error(fmt.Sprintf(format, req.Method, req.URL.String(), err, code, msg))
}

func (h *handler) logRequest(req *http.Request, code int, err error) {
	msg := http.StatusText(code)
	if err != nil {
		const format = "%s %s    <-    %d %s    error=%s"
		h.Log.Error(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg, err))
		return
	}
	const format = "%s %s    <-    %d %s"
	h.Log.Info(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg))
}
