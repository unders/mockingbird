package html

import (
	"fmt"

	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// Template implements the generation of HTML pages
type Template struct {
	TemplateDir string
}

//
// Success pages
//

// Dashboard returns Dashboard page
func (t *Template) Dashboard(d mockingbird.Dashboard) ([]byte, error) {
	return []byte(fmt.Sprintf("Dashbord page: %+v", d)), nil
}

// ListTests returns test results page
func (t *Template) ListTest(r *mockingbird.TestResults) ([]byte, error) {
	return []byte(fmt.Sprintf("Test results  page: %+v", r)), nil
}

// ShowTest returns test result page
func (t *Template) ShowTest(r mockingbird.TestResult) ([]byte, error) {
	return []byte(fmt.Sprintf("Test result  page: %+v", r)), nil
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
