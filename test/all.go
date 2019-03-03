// +build mage

package main

import (
	"github.com/magefile/mage/mg"
)

type All mg.Namespace

// Test tests all services
func (All) Test() {
	mg.SerialDeps(
		Google.Test,
	)
}
