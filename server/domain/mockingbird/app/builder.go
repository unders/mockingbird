package app

import (
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/unders/mockingbird/server/pkg/handler"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/domain/mockingbird/html"
	"github.com/unders/mockingbird/server/domain/mockingbird/mock"
)

// Options defines the required input to function app.Create
type Options struct {
	Env         mockingbird.Env
	Logger      *log.Logger
	FaviconDir  string
	TemplateDir string
	AssetDir    string
}

// Create creates the application
func Create(o Options) (*Builder, error) {
	tmpl, err := newTemplate(o.Env, o.TemplateDir)
	if err != nil {
		return nil, err
	}

	b := Builder{
		log: o.Logger,
		// TODO: Change to real app when it is implemented
		// app:     &Mockingbird{},
		app:     &mock.AppMockingbird{Now: time.Now().UTC()},
		favicon: handler.Favicons(o.FaviconDir),
		assets:  http.Dir(o.AssetDir),
		tmpl:    tmpl,
	}
	return &b, nil
}

// Builder builds the application
type Builder struct {
	favicon func(*http.Request) (http.Handler, bool)
	app     mockingbird.App
	log     *log.Logger
	tmpl    *html.Template
	assets  http.FileSystem
}

// Favicon returns a http.Handler for favicons
func (b *Builder) Favicon() func(*http.Request) (http.Handler, bool) {
	return b.favicon
}

// Assets returns a http.FileSystem for all public assets
func (b *Builder) Assets() http.FileSystem {
	return b.assets
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

//
// private
//
func newTemplate(env mockingbird.Env, templateDir string) (tmpl *html.Template, err error) {
	if mockingbird.DEV == env {
		tmpl, err = html.NewReloadableTemplate(templateDir)
		err = errors.Wrapf(err, "html.NewReloadableTemplate(%s) failed", templateDir)
		return tmpl, err
	}

	tmpl, err = html.NewTemplate(templateDir)
	err = errors.Wrapf(err, "html.NewsTemplate(%s) failed", templateDir)
	return tmpl, err
}
