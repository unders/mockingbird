// +build mage

// Tests the github.com/unders/mockingbird code.
package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Server mg.Namespace

// All runs all server targets
func (Server) All() {
	mg.SerialDeps(
		Go.Version,
		Server.download,
		Server.Test,
		Server.Race,
	)
}

// Test runs: go test github.com/unders/mockingbird/cmd/...
func (Server) Test() error {
	mg.SerialDeps(Server.download)
	return sh.RunV("go", "test", "./server/...")
}

// Race runs: go test github.com/unders/mockingbird/cmd/...
func (Server) Race() error {
	mg.SerialDeps(Server.download)
	return sh.RunV("go", "test", "-race", "./server/...")
}

//
// Private commands
//
func (Server) download() error {
	mg.SerialDeps(Server.root)
	return sh.Run("go", "mod", "download")
}

func (Server) root() error {
	return os.Chdir("../")
}
