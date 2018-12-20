package testdata

import "flag"

// Update is true if something should be updated in a package
//
// Note:
//
//          Only activate this flag when testing a specific package;
//
// Usage:
//          go test github.com/unders/mockingbird/pkg/html -update
//
var Update = flag.Bool("update", false, "update golden files")

func init() {
	flag.Parse()
}
