package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/unders/mockingbird/cmd/domain/mockinbird"
	"github.com/unders/mockingbird/cmd/pkg/rest"
)

type Handler struct {
	HTML mockinbird.HTMLAdapter
	Log  mockinbird.Log
}

// GET  http://localhost:8080/ -> redirect-to: http://localhost:8080/v1/dashboard
// GET  http://localhost:8080/v1/dashboard
// POST http://localhost:8080/v1/environment/{env}/tests/
// GET  http://localhost:8080/v1/environment/{env}/tests/{ID}
// GET  http://localhost:8080/v1/environment/{env}/tests/
// GET  http://localhost:8080/v1/environment/{env}/tests/?service=sim-management

// POST http://localhost:8080/v1/environment/{env}/tests/-/services/sim-management
func (h Handler) Make() http.Handler {
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
		case rest.Route{Method: http.MethodPost, Path: "/v1/environment/*/tests"}:
			h.runTest(w, req, path)
		case rest.Route{Method: http.MethodGet, Path: "/v1/environment/*/tests/*"}:
			h.showTestResult(w, req, path)
		case rest.Route{Method: http.MethodGet, Path: "/v1/environment/*/tests"}:
			h.listTestResults(w, req, path)
		case rest.Route{Method: http.MethodPost, Path: "/v1/environment/*/tests/-/services/*"}:
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
func (h *Handler) showDashboard(w http.ResponseWriter, req *http.Request) {
	code, b, err := h.HTML.Dashboard()
	h.write(w, req, code, b, err)
}

func (h *Handler) runTest(w http.ResponseWriter, req *http.Request, path rest.Path) {
	env := path.String(2, "")
	if err := h.HTML.HasEnvError(env); err != nil {
		h.invalidEnvironment(w, req, env, err)
		return
	}

	id, err := h.HTML.RunTest(env)
	if err != nil {
		h.internalError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/v1/environment/%s/tests/%d", env, id), http.StatusSeeOther)
}

func (h *Handler) showTestResult(w http.ResponseWriter, req *http.Request, path rest.Path) {
	env := path.String(2, "")
	if err := h.HTML.HasEnvError(env); err != nil {
		h.invalidEnvironment(w, req, env, err)
		return
	}

	id, err := path.Int(4)
	if err != nil {
		h.invalidID(w, req, id, err)
		return
	}

	code, b, err := h.HTML.ShowTestResult(env, id)
	h.write(w, req, code, b, err)
}

func (h *Handler) listTestResults(w http.ResponseWriter, req *http.Request, path rest.Path) {
	env := path.String(2, "")
	if err := h.HTML.HasEnvError(env); err != nil {
		h.invalidEnvironment(w, req, env, err)
		return
	}
	code, b, err := h.HTML.ListTestResults(env)
	h.write(w, req, code, b, err)
}

func (h *Handler) runTestForService(w http.ResponseWriter, req *http.Request, path rest.Path) {
	env := path.String(2, "")
	if err := h.HTML.HasEnvError(env); err != nil {
		h.invalidEnvironment(w, req, env, err)
		return
	}

	service := path.String(6, "")
	if err := h.HTML.HasServiceError(env); err != nil {
		h.invalidEnvironment(w, req, env, err)
		return
	}
	id, err := h.HTML.RunTestForService(env, service)
	if err != nil {
		h.internalError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/v1/environment/%s/tests/%d", env, id), http.StatusSeeOther)
}

//
// HTML Errors
//

func (h *Handler) internalError(w http.ResponseWriter, req *http.Request, err error) {
	code, b := h.HTML.ErrorServer()
	h.write(w, req, code, b, err)
}

func (h *Handler) invalidID(w http.ResponseWriter, req *http.Request, id int, err error) {
	_, b := h.HTML.ErrorClient(
		"Not a Number",
		fmt.Sprintf("ID %d is not a valid number", id),
	)

	h.write(w, req, http.StatusNotFound, b, err)
}

func (h *Handler) invalidEnvironment(w http.ResponseWriter, req *http.Request, env string, err error) {
	_, b := h.HTML.ErrorClient(
		"Not Found",
		fmt.Sprintf("The environment %s does not exist.", env),
	)

	h.write(w, req, http.StatusNotFound, b, err)
}

func (h *Handler) invalidService(w http.ResponseWriter, req *http.Request, service string, err error) {
	_, b := h.HTML.ErrorClient(
		"Not Found",
		fmt.Sprintf("The service %s does not exist.", service),
	)

	h.write(w, req, http.StatusNotFound, b, err)
}

func (h *Handler) invalidURLFormat(w http.ResponseWriter, req *http.Request, err error) {
	_, b := h.HTML.ErrorClient(
		"Invalid URL",
		"The URL has invalid characters.",
	)

	h.write(w, req, http.StatusBadRequest, b, err)
}

func (h *Handler) routeNotFound(w http.ResponseWriter, req *http.Request) {
	_, b := h.HTML.ErrorNotFound()
	h.write(w, req, http.StatusBadRequest, b, errors.New("route not found"))
}

//
// HTML Writer
//

func (h *Handler) write(w http.ResponseWriter, req *http.Request, code int, buf []byte, err error) {
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

func (h *Handler) isJSON(req *http.Request) bool {
	ct := req.Header.Get("Content-Type")
	return ct == "application/json"
}

func (h *Handler) jsonNotImplemented(w http.ResponseWriter, req *http.Request) {
	const code = http.StatusNotImplemented
	var status = http.StatusText(code)
	const msg = "JSON API not implemented by the server."

	const format = `{"error": { "code": %d, "status": "%s", "message": "%s" } }`
	b := []byte(fmt.Sprintf(format, code, status, msg))
	err := errors.Errorf(format, code, status, msg)
	h.writeJSON(w, req, code, b, err)
}

func (h *Handler) writeJSON(w http.ResponseWriter, req *http.Request, code int, buf []byte, err error) {
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
func (h *Handler) logResponseFailure(req *http.Request, code int, err error) {
	msg := http.StatusText(code)
	const format = "%s %s    response error=%s    [%d %s]"
	h.Log.Error(fmt.Sprintf(format, req.Method, req.URL.String(), err, code, msg))
}

func (h *Handler) logRequest(req *http.Request, code int, err error) {
	msg := http.StatusText(code)
	if err != nil {
		const format = "%s %s    <-    %d %s    error=%s"
		h.Log.Error(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg, err))
		return
	}
	const format = "%s %s    <-    %d %s"
	h.Log.Info(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg))
}
