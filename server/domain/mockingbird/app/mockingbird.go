package app

import (
	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// Mockingbird implements the mockingbird.App interface
//
// Note:
//
//        app.Mockingbird implements the business logic of the application
//
//
type Mockingbird struct {
}

// Verifies that *Mockingbird implements mockingbird.Builder interface
var _ mockingbird.App = &Mockingbird{}

func (m *Mockingbird) Dashboard() error {
	return nil
}

func (m *Mockingbird) RunTest() (id string, err error) {
	return "", nil
}
func (m *Mockingbird) ShowTestResult(id string) error {
	return nil
}
func (m *Mockingbird) ListTestResults(service string) (err error) {
	return nil
}
func (m *Mockingbird) RunTestForService(service string) (id string, err error) {
	return "", nil
}
