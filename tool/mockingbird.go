// +build mage

package main

import (
	"os"

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

	if err := os.Chdir("test"); err != nil {
		return err
	}

	return sh.RunV("modd", "-f", "../tool/app/mockingbird/dev.conf")
}
