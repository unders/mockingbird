// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Mockingbird defines mockingbird specific targets
type Mockingbird mg.Namespace

// Start starts mockingbird application
func (Mockingbird) Start() error {
	mg.SerialDeps(
		Server.Build,
	)

	return sh.RunV("modd", "-f", "./tool/app/mockingbird/dev.conf")
}
