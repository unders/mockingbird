package html

import "github.com/unders/mockingbird/cmd/domain/mockinbird"

// Adapter implements the mockingbird.HTMLAdapter interface
//
//
// Note:
//
//         * It adapts between the http Handler and the mockingbird.Container
//         * This adapter returns HTML pages.
//
//
type Adapter struct {
	container mockinbird.Container
}

// Verifies that *Adapter implements mockingbird.HTMLAdapter interface
// var _ mockingbird.HTMLAdapter = &Adapter{}
