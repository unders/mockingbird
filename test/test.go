// +build mage

package main

import (
	"github.com/magefile/mage/mg"
)

type Test mg.Namespace

// All tests all services
func (Test) All() {
	mg.SerialDeps(
		Google.Test,
	)
}
