package handler

import (
	"net/http"
	"path/filepath"
)

// Favicons returns a function that serves favicons from the specified directory when called.
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

	h := map[string]http.Handler{
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

	return func(req *http.Request) (http.Handler, bool) {
		handler, found := h[req.URL.EscapedPath()]
		return handler, found
	}
}
