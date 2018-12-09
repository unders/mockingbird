package html

import "github.com/unders/mockingbird/server/domain/mockingbird"

// Adapter implements the mockingbird.HTMLAdapter interface
//
//
// Note:
//
//         * It adapts between the http.Handler and the mockingbird.App
//         * This adapter returns HTML pages.
//
//
type Adapter struct {
	container mockingbird.App
}

// Verifies that *Adapter implements mockingbird.HTMLAdapter interface
// var _ mockingbird.HTMLAdapter = &Adapter{}
