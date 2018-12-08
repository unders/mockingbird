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
		// All other Content Types are treated as HTML.

		path, route, err := router.New(req)
		if err != nil {
			h.invalidURLFormat(w, req, err)
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
			h.routeNotFound(w, req)
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
	id, err := h.HTML.RunTest()
	if err != nil {
		h.internalError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/v1/tests/%d", id), http.StatusSeeOther)
}

func (h *handler) showTestResult(w http.ResponseWriter, req *http.Request, path rest.Path) {
	id, err := path.Int(2)
	if err != nil {
		h.invalidID(w, req, id, err)
		return
	}

	code, b, err := h.HTML.ShowTestResult(id)
	h.write(w, req, code, b, err)
}

func (h *handler) listTestResults(w http.ResponseWriter, req *http.Request) {
	code, b, err := h.HTML.ListTestResults()
	h.write(w, req, code, b, err)
}

func (h *handler) runTestForService(w http.ResponseWriter, req *http.Request, path rest.Path) {
	service := path.String(4, "")
	if err := h.HTML.HasServiceError(service); err != nil {
		h.invalidService(w, req, service, err)
		return
	}
	id, err := h.HTML.RunTestForService(service)
	if err != nil {
		h.internalError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/v1/tests/%d", id), http.StatusSeeOther)
}

//
// HTML Errors
//

func (h *handler) internalError(w http.ResponseWriter, req *http.Request, err error) {
	code, b := h.HTML.ErrorServer()
	h.write(w, req, code, b, err)
}

func (h *handler) invalidID(w http.ResponseWriter, req *http.Request, id int, err error) {
	_, b := h.HTML.ErrorClient(
		"Not a Number",
		fmt.Sprintf("ID %d is not a valid number", id),
	)

	h.write(w, req, http.StatusBadRequest, b, err)
}

func (h *handler) invalidService(w http.ResponseWriter, req *http.Request, service string, err error) {
	_, b := h.HTML.ErrorClient(
		"Not Found",
		fmt.Sprintf("The service %s does not exist.", service),
	)

	h.write(w, req, http.StatusNotFound, b, err)
}

func (h *handler) invalidURLFormat(w http.ResponseWriter, req *http.Request, err error) {
	_, b := h.HTML.ErrorClient(
		"Invalid URL",
		"The URL has invalid characters.",
	)

	h.write(w, req, http.StatusBadRequest, b, err)
}

func (h *handler) routeNotFound(w http.ResponseWriter, req *http.Request) {
	_, b := h.HTML.ErrorNotFound()
	h.write(w, req, http.StatusBadRequest, b, errors.New("route not found"))
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
