package html

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/unders/mockingbird/server/pkg/html"

	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// Layouts
const (
	mainLayout = "main.html"
)

// Page file names
const (
	dashboard string = "dashboard.html"
)

//
// Path
//

// Path defines the applications URL
type Path struct {
	Dashboard      string
	ListTests      string
	showTest       string
	ListTestSuites string
}

var path = Path{
	Dashboard:      "/dashboard",
	ListTests:      "/tests/",
	showTest:       "/tests/%s",
	ListTestSuites: "/tests/-/suites/",
}

// ShowTest returns path to show test page
func (p Path) ShowTest(id string) string {
	return fmt.Sprintf(p.showTest, id)
}

//
// Pages
//

// dashboardPage contains all required data for rendering dashboard page
type dashboardPage struct {
	Title string
	Path  *Path
	mockingbird.Dashboard
}

// NewReloadableTemplate returns html.Template
func NewReloadableTemplate(templateDir string) (*Template, error) {
	tmpl, err := html.NewReloadableTemplate(templateDir)
	if err != nil {
		return nil, errors.Wrapf(err, "html.NewReloadableTemplate(%s) failed", templateDir)
	}
	return &Template{tmpl: tmpl}, nil
}

// NewTemplate returns html.Template
func NewTemplate(templateDir string) (*Template, error) {
	tmpl, err := html.NewTemplate(templateDir)
	if err != nil {
		return nil, errors.Wrapf(err, "html.NewTemplate(%s) failed", templateDir)
	}
	return &Template{tmpl: tmpl}, nil
}

// Template implements the generation of HTML pages
type Template struct {
	tmpl html.Templater
}

//
// Success pages
//

// Dashboard returns the dashboard page as HTML
func (t *Template) Dashboard(d mockingbird.Dashboard) ([]byte, error) {
	const title = "Dashboard - Mockingbird"
	page := dashboardPage{Title: title, Dashboard: d, Path: &path}
	return t.tmpl.Execute(mainLayout, dashboard, page)
}

// ListTests returns test results page
func (t *Template) ListTest(r *mockingbird.TestResults) ([]byte, error) {
	return []byte(fmt.Sprintf("Test results  page: %+v", r)), nil
}

// ShowTest returns test result page
func (t *Template) ShowTest(r mockingbird.TestResult) ([]byte, error) {
	return []byte(fmt.Sprintf("Test result  page: %+v", r)), nil
}

// ShowTestSuites returns test  suite page
func (t *Template) ShowTestSuites(ts []mockingbird.TestSuite) ([]byte, error) {
	return []byte(fmt.Sprintf("Test suites page: %+v", ts)), nil
}

//
// Client Errors
//

// ErrorNotFound returns error not found page
func (t *Template) ErrorNotFound() []byte {
	return []byte("404 Not Found")
}

// InvalidURL returns an invalid URL page
func (t *Template) InvalidURL() []byte {
	return []byte("Invalid URL")
}

//
// Server Error
//

// InternalError returns Internal Error Page
func (t *Template) InternalError() []byte {
	return []byte("Internal Error Page")
}
