package html

import (
	"fmt"
	"time"

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
	dashboard       = "dashboard.html"
	testResult      = "tests/show.html"
	testResultsFile = "tests/index.html"
	testSuitesFile  = "tests/suites.html"
	errorFile       = "error.html"
)

// Assets files
const (
	cssFile = "/public/css/main.css"
)

//
// Path
//

// Path defines the applications URL
type Path struct {
	Dashboard      string
	ListTests      string
	RunTest        string
	showTest       string
	ListTestSuites string
}

var path = Path{
	Dashboard:      "/dashboard",
	ListTests:      "/tests/",
	RunTest:        "/tests/",
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
	CSS   string
	Title string

	PageTitle string

	FullTestSuiteState mockingbird.TestResult
	TestSuiteState     mockingbird.TestResult

	ReloadPath              string
	LatestTestSuitePath     string
	LatestFullTestSuitePath string

	Path *Path
	mockingbird.Dashboard
}

// testResultPage contains all required data for rendering test result page
type testResultPage struct {
	CSS             string
	Title           string
	PageTitle       string
	ReloadPath      string
	Path            *Path
	Result          mockingbird.TestResult
	ResultLog       string
	ResultStartTime string
}

type testSuites struct {
	CSS        string
	Title      string
	PageTitle  string
	ReloadPath string
	Path       *Path

	TestSuites []mockingbird.TestSuite
}

type testResultsPage struct {
	CSS        string
	Title      string
	PageTitle  string
	ReloadPath string
	Path       *Path

	MoreResult bool

	NextPage    string
	TestResults []mockingbird.TestResult
}

type errorPage struct {
	CSS         string
	Title       string
	PageTitle   string
	Description string

	Path *Path
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

	page := dashboardPage{
		Title: title,
		CSS:   cssFile,

		Path: &path,

		PageTitle:               "Dashboard",
		ReloadPath:              path.Dashboard,
		LatestTestSuitePath:     path.ShowTest(string(d.Stats.LatestDoneTestSuiteID)),
		LatestFullTestSuitePath: path.ShowTest(string(d.Stats.LatestDoneFullTestSuiteID)),
		FullTestSuiteState:      mockingbird.TestResult{State: d.Stats.LatestDoneFullTestSuiteState},
		TestSuiteState:          mockingbird.TestResult{State: d.Stats.LatestDoneTestSuiteState},

		Dashboard: d,
	}
	return t.tmpl.Execute(mainLayout, dashboard, page)
}

// ListTests returns test results page
func (t *Template) ListTest(ts *mockingbird.TestResults) ([]byte, error) {
	const title = "Test history - Mockingbird"

	page := testResultsPage{
		Title: title,
		CSS:   cssFile,

		ReloadPath: path.ListTests,
		PageTitle:  "Test History",
		Path:       &path,

		MoreResult: len(ts.NextPageToken) != 0,

		NextPage:    fmt.Sprintf("%s?page_token=%s", path.ListTests, ts.NextPageToken),
		TestResults: ts.TestResults,
	}
	return t.tmpl.Execute(mainLayout, testResultsFile, page)
}

// ShowTest returns test result page
func (t *Template) ShowTest(ts mockingbird.TestResult) ([]byte, error) {
	const title = "Test result - Mockingbird"

	page := testResultPage{
		Title: title,
		CSS:   cssFile,

		ReloadPath:      path.ShowTest(string(ts.ID)),
		PageTitle:       fmt.Sprintf("%s", ts.TestSuite),
		Path:            &path,
		Result:          ts,
		ResultLog:       string(ts.Log),
		ResultStartTime: ts.StartTime.Format(time.RFC3339),
	}
	return t.tmpl.Execute(mainLayout, testResult, page)
}

// ShowTestSuites returns test  suite page
func (t *Template) ShowTestSuites(ts []mockingbird.TestSuite) ([]byte, error) {
	const title = "Test suites - Mockingbird"

	page := testSuites{
		Title: title,
		CSS:   cssFile,

		ReloadPath: path.ListTestSuites,
		PageTitle:  "Test Suites",
		Path:       &path,
		TestSuites: ts,
	}
	return t.tmpl.Execute(mainLayout, testSuitesFile, page)
}

//
// Client Errors
//

// ErrorNotFound returns error not found page
func (t *Template) ErrorNotFound() []byte {
	const title = "Page Not Found - Mockingbird"

	page := errorPage{
		Title: title,
		CSS:   cssFile,
		Path:  &path,

		PageTitle:   "Page Not Found",
		Description: "The page might be gone due to a new deployment :).",
	}
	buf, err := t.tmpl.Execute(mainLayout, errorFile, page)
	if err != nil {
		return internalError(err)
	}

	return buf
}

// InvalidURL returns an invalid URL page
func (t *Template) InvalidURL() []byte {
	const title = "Invalid URL - Mockingbird"

	page := errorPage{
		Title: title,
		CSS:   cssFile,
		Path:  &path,

		PageTitle:   "Invalid URL",
		Description: "You wrote an invalid URL, please try again...",
	}
	buf, err := t.tmpl.Execute(mainLayout, errorFile, page)
	if err != nil {
		return internalError(err)
	}

	return buf
}

//
// Server Error
//

// InternalError returns Internal Error Page
func (t *Template) InternalError() []byte {
	const title = "Internal Error - Mockingbird"

	page := errorPage{
		Title: title,
		CSS:   cssFile,
		Path:  &path,

		PageTitle:   "Internal Error",
		Description: "The server had an internal error.",
	}
	buf, err := t.tmpl.Execute(mainLayout, errorFile, page)
	if err != nil {
		return internalError(err)
	}
	return buf
}

func internalError(err error) []byte {
	s := fmt.Sprintf("Internal Error %s", err)
	return []byte(s)
}
