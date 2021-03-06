package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/unders/mockingbird/server/domain/mockingbird"

	"github.com/pkg/errors"
	"github.com/unders/mockingbird/server/pkg/rest"
)

type handler struct {
	Favicon       func(r *http.Request) (http.Handler, bool)
	notAuthorized func(w http.ResponseWriter, r *http.Request) bool
	Assets        http.FileSystem
	AssetsPrefix  string
	HTML          mockingbird.HTMLAdapter
	Log           mockingbird.Log
}

func createHandler(h handler) http.Handler {
	router := rest.Router{}

	assets := "/public/"
	if h.AssetsPrefix != "" {
		assets = h.AssetsPrefix
	}

	if h.notAuthorized == nil {
		h.notAuthorized = notAuthorized
	}

	fs := http.StripPrefix(assets, http.FileServer(h.Assets))

	f := func(w http.ResponseWriter, req *http.Request) {
		if h.isJSON(req) {
			h.jsonNotImplemented(w, req)
			return
		}
		if strings.HasPrefix(req.URL.Path, assets) {
			fs.ServeHTTP(w, req)
			return
		}

		if h.notAuthorized(w, req) {
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
		//
		// GET
		//
		case rest.Route{Method: http.MethodGet, Path: ""}:
			url := "/dashboard"
			http.Redirect(w, req, url, http.StatusSeeOther)
		case rest.Route{Method: http.MethodGet, Path: "/dashboard"}:
			h.showDashboard(w, req)
		case rest.Route{Method: http.MethodGet, Path: "/tests"}:
			h.listTests(w, req)
		case rest.Route{Method: http.MethodGet, Path: "/tests/*"}:
			h.showTest(w, req, path)
		case rest.Route{Method: http.MethodGet, Path: "/tests/*/suites"}:
			h.showTestSuites(w, req, path)
			//
			// POST
			//
		case rest.Route{Method: http.MethodPost, Path: "/tests"}:
			h.runTest(w, req)
		default:
			if favicon, found := h.Favicon(req); found {
				h.logRequest(req, http.StatusOK, nil)
				favicon.ServeHTTP(w, req)
				return
			}

			err := errors.New("route not found")
			h.write(w, req, http.StatusNotFound, h.HTML.ErrorNotFound(), err)
		}
	}

	return http.HandlerFunc(f)
}

//
// HTML Handlers
//

func notAuthorized(w http.ResponseWriter, r *http.Request) bool {
	user, pass, _ := r.BasicAuth()

	if "user" != user || "password" != pass {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return true
	}
	return false
}

//
func (h *handler) showDashboard(w http.ResponseWriter, req *http.Request) {
	code, b, err := h.HTML.Dashboard()
	h.write(w, req, code, b, err)
}

func (h *handler) listTests(w http.ResponseWriter, req *http.Request) {
	code, b, err := h.HTML.ListTests(req.URL.Query().Get("page_token"))
	h.write(w, req, code, b, err)
}

func (h *handler) showTest(w http.ResponseWriter, req *http.Request, path rest.Path) {
	id := path.String(1, "")
	code, b, err := h.HTML.ShowTest(mockingbird.ULID(id))

	h.write(w, req, code, b, err)
}

func (h *handler) showTestSuites(w http.ResponseWriter, req *http.Request, path rest.Path) {
	code, b, err := h.HTML.ShowTestSuites()
	h.write(w, req, code, b, err)
}

func (h *handler) runTest(w http.ResponseWriter, req *http.Request) {
	ts := req.FormValue("test_suite")
	id, code, body, err := h.HTML.RunTest(mockingbird.TestSuite(ts))
	if err != nil {
		h.write(w, req, code, body, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/tests/%s", id), http.StatusSeeOther)
}

//
// HTML Writer
//

func (h *handler) write(w http.ResponseWriter, req *http.Request, code int, buf []byte, err error) {
	body := bytes.NewReader(buf)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	w.WriteHeader(code)
	if _, err := io.Copy(w, body); err != nil {
		err = errors.Wrap(err, "io.Copy(w, body) failed")
		h.logResponseFailure(req, code, err)
		return
	}

	h.logRequest(req, code, err)
}

//
// JSON
//

func (h *handler) isJSON(req *http.Request) bool {
	const jsonFormat = "application/json"
	accept := req.Header.Get("Accept")
	ct := req.Header.Get("Content-Type")

	return ct == jsonFormat || accept == jsonFormat
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	w.WriteHeader(code)
	if _, err := io.Copy(w, body); err != nil {
		err = errors.Wrap(err, "io.Copy(w, body) failed")
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
		if code < http.StatusInternalServerError {
			const format = "%s %s    <-    %d %s    notice=%s"
			h.Log.Notice(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg, err))
			return
		}

		const format = "%s %s    <-    %d %s    error=%s"
		h.Log.Error(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg, err))
		return
	}
	const format = "%s %s    <-    %d %s"
	h.Log.Info(fmt.Sprintf(format, req.Method, req.URL.String(), code, msg))
}
