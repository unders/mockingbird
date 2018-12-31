package app

import (
	"log"
	"net/http"
	"time"

	"github.com/unders/mockingbird/server/pkg/handler"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/domain/mockingbird/html"
	"github.com/unders/mockingbird/server/domain/mockingbird/mock"
)

// Options defines the required input to function app.Create
type Options struct {
	Logger      *log.Logger
	FaviconDir  string
	TemplateDir string
}

// Create creates the application
func Create(o Options) (*Builder, error) {
	b := Builder{
		log: o.Logger,
		// TODO: Change to real app when it is implemented
		// app:     &Mockingbird{},
		app:     &mock.AppMockingbird{Now: time.Now().UTC()},
		favicon: handler.Favicons(o.FaviconDir),
		tmpl:    &html.Template{TemplateDir: o.TemplateDir},
	}
	return &b, nil
}

// Builder builds the application
type Builder struct {
	favicon func(*http.Request) (http.Handler, bool)
	app     mockingbird.App
	log     *log.Logger
	tmpl    *html.Template
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
	return &html.Adapter{App: b.app, Tmpl: b.tmpl}
}

// Log returns the mockingbird.Log
func (b *Builder) Log() mockingbird.Log {
	return &mockingbird.Logger{Log: b.log}
}
