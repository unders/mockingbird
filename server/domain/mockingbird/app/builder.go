package app

import (
	"log"
	"net/http"

	"github.com/unders/mockingbird/server/pkg/handler"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/domain/mockingbird/html"
	"github.com/unders/mockingbird/server/domain/mockingbird/mock"
)

// Create creates the application
func Create(b Builder) (*Builder, error) {
	// do the setup here
	b.app = &Mockingbird{}
	b.favicon = handler.Favicons(b.FaviconDir)
	return &b, nil
}

// Builder builds the application
type Builder struct {
	// Required input
	Logger     *log.Logger
	FaviconDir string

	// constructed in create method
	favicon func(*http.Request) (http.Handler, bool)
	app     mockingbird.App
}

// Favicon returns a http.Handler for favicons
func (b *Builder) Favicon() func(*http.Request) (http.Handler, bool) {
	return b.favicon
}

// HTMLMockAdapter returns a mock.HTMLAdapter
func (b *Builder) HTMLMockAdapter() mockingbird.HTMLAdapter {
	return mock.HTMLAdapter{Code: 200, Body: []byte("Hello World!")}
}

// HTMLAdapter returns the mockingbird.HTMLAdapter
func (b *Builder) HTMLAdapter() mockingbird.HTMLAdapter {
	return &html.Adapter{App: b.app}
}

// Log returns the mockingbird.Log
func (b *Builder) Log() mockingbird.Log {
	return &mockingbird.Logger{Log: b.Logger}
}
