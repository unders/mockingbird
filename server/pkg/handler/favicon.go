package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
)

// FaviconBox returns a function that serves favicons from the specified directory.
//
// Usage:
//         box := packr.New("Favicons", "web/mockingbird/public/favicon")
//         f := handler.Favicons(box)
//         if handler, found := f(req); found {
//             handler.ServeHTTP(w, req)
//         }
//
func FaviconsBox(box *packr.Box) func(r *http.Request) (http.Handler, bool) {
	panic("handler.FaviconBox() not working at the moment; waiting for a stable v2 release.")
	serveFile := func(filename string) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			// A 2 day cache! (if the favicon is updated it will take
			// 2 day before a user sees the change)
			w.Header().Set("Cache-Control", "public, max-age=172800")
			http.FileServer(box).ServeHTTP(w, r)
		}

		return http.HandlerFunc(h)
	}
	mux := faviconMux(serveFile)

	return func(req *http.Request) (http.Handler, bool) {
		if http.MethodGet != req.Method {
			return nil, false
		}
		handler, found := mux[req.URL.EscapedPath()]
		return handler, found
	}
}

// Favicons returns a function that serves favicons from the specified directory.
//
// Usage:
//         f := handler.Favicons("web/mockingbird/public/favicon")
//         if handler, found := f(req); found {
//             handler.ServeHTTP(w, req)
//         }
//
func Favicons(dir string) func(r *http.Request) (http.Handler, bool) {

	serveFile := func(filename string) http.Handler {
		file := filepath.Join(dir, filename)

		h := func(w http.ResponseWriter, r *http.Request) {
			// A 2 day cache! (if the favicon is updated it will take
			// 2 day before a user sees the change)
			w.Header().Set("Cache-Control", "public, max-age=172800")
			http.ServeFile(w, r, file)
		}

		return http.HandlerFunc(h)
	}

	h := faviconMux(serveFile)

	return func(req *http.Request) (http.Handler, bool) {
		if http.MethodGet != req.Method {
			return nil, false
		}

		handler, found := h[req.URL.EscapedPath()]
		return handler, found
	}
}

func faviconMux(serveFile func(filename string) http.Handler) map[string]http.Handler {
	return map[string]http.Handler{
		"/android-chrome-192x192.png":       serveFile("android-chrome-192x192.png"),
		"/android-chrome-512x512.png":       serveFile("android-chrome-512x512.png"),
		"/apple-touch-icon.png":             serveFile("apple-touch-icon.png"),
		"/apple-touch-icon-precomposed.png": serveFile("apple-icon-precomposed.png"),
		"/apple-icon-precomposed.png":       serveFile("apple-icon-precomposed.png"),
		"/browserconfig.xml":                serveFile("browserconfig.xml"),
		"/favicon.ico":                      serveFile("favicon.ico"),
		"/favicon-16x16.png":                serveFile("favicon-16x16.png"),
		"/favicon-32x32.png":                serveFile("favicon-32x32.png"),
		"/mstile-70x70.png":                 serveFile("mstile-70x70.png"),
		"/mstile-144x144.png":               serveFile("mstile-144x144.png"),
		"/mstile-150x150.png":               serveFile("mstile-150x150.png"),
		"/mstile-310x150.png":               serveFile("mstile-310x150.png"),
		"/mstile-310x310.png":               serveFile("mstile-310x310.png"),
		"/safari-pinned-tab.svg":            serveFile("safari-pinned-tab.svg"),
		"/site.webmanifest.json":            serveFile("site.webmanifest.json"),
	}
}
