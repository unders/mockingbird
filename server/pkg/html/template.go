package html

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/unders/mockingbird/server/pkg/errs"

	"bytes"

	"fmt"

	"github.com/pkg/errors"
)

// Templater defines the template interface
type Templater interface {
	Execute(layout, filepath string, data interface{}) ([]byte, error)
}

// Template assembles a layout, parts, a page and data into a HTML page.
type Template struct {
	Template *template.Template
	Page     map[string]string
}

// NewTemplate returns an html.Templater
func NewTemplate(root string) (Templater, error) {
	return newTemplate(root)
}

// newTemplate returns a *Template
func newTemplate(root string) (*Template, error) {
	var (
		err error
		t   *template.Template
	)

	l := filepath.Join(root, "layout", "*")
	if t, err = template.ParseGlob(l); err != nil {
		return nil, errors.WithStack(err)
	}

	var parts []string
	partdir := filepath.Join(root, "part")
	walkparts := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "parts file walk failed")
		}
		if f.IsDir() {
			return nil
		}

		parts = append(parts, path)
		return nil
	}

	if err = filepath.Walk(partdir, walkparts); err != nil {
		return nil, errors.WithStack(err)
	}

	if t, err = t.ParseFiles(parts...); err != nil {
		return nil, errors.WithStack(err)
	}

	page := map[string]string{}
	pagedir := filepath.Join(root, "page")
	walk := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "file walk failed")
		}
		if f.IsDir() {
			return nil
		}

		key := path[len(pagedir)+1:]
		p, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "ReadFile failed")
		}
		page[key] = string(p)

		return nil
	}
	if err = filepath.Walk(pagedir, walk); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Template{Template: t, Page: page}, nil
}

// Execute parses a file with the given layout template
//
// Usage:
//
//        data := map[string]string{ "Title": "A Title" }
//        t := html.NewTemplate("resource/template/")
//        b, err := t.Execute("main.html", "post.html", data)
//
func (t *Template) Execute(layout, filepath string, data interface{}) ([]byte, error) {
	tmpl := t.Template.Lookup(layout)
	if tmpl == nil {
		msg := fmt.Sprintf("layout %s not found", layout)
		return nil, errs.NotFound(msg)
	}

	page, ok := t.Page[filepath]
	if !ok {
		msg := fmt.Sprintf("page %s not found", filepath)
		return nil, errs.NotFound(msg)
	}

	tmpl, err := tmpl.Clone()
	if err != nil {
		return nil, errors.Errorf("Could not clone %s", layout)
	}

	tmpl, err = tmpl.Parse(page)
	if err != nil {
		msg := fmt.Sprintf("template %s page %s could not parse content %s", layout, filepath, page)
		return nil, errors.Wrap(err, msg)
	}

	// Execution stops immediately with an error.
	tmpl = tmpl.Option("missingkey=error")

	b := &bytes.Buffer{}
	if err := tmpl.Execute(b, data); err != nil {
		return nil, errors.WithStack(err)
	}

	return b.Bytes(), nil
}

//
// ReloadableTemplate
//

// NewReloadableTemplate returns an html.Templater
func NewReloadableTemplate(root string) (Templater, error) {
	return &ReloadableTemplate{root: root}, nil
}

// ReloadableTemplate is used to generate a HTML page
type ReloadableTemplate struct {
	root string
}

// Execute parses a file with the given layout template
func (t *ReloadableTemplate) Execute(layout, filepath string, data interface{}) ([]byte, error) {
	tmpl, err := newTemplate(t.root)
	if err != nil {
		return nil, err
	}
	return tmpl.Execute(layout, filepath, data)
}
