// +build mage

// Tests the github.com/unders/mockingbird code.
package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace
type Cmd mg.Namespace

// Cmd runs  all CMD targets
func (Test) Cmd() {
	mg.SerialDeps(
		Go.Version,
		Cmd.download,
		Cmd.Test,
		Cmd.Race,
	)
}

// Test runs: go test github.com/unders/mockingbird/cmd/...
func (Cmd) Test() error {
	mg.SerialDeps(Cmd.download)
	return sh.RunV("go", "test", "./cmd/...")
}

// Race runs: go test github.com/unders/mockingbird/cmd/...
func (Cmd) Race() error {
	mg.SerialDeps(Cmd.download)
	return sh.RunV("go", "test", "-race", "./cmd/...")
}

//
// Private commands
//
func (Cmd) download() error {
	mg.SerialDeps(Cmd.root)
	return sh.Run("go", "mod", "download")
}

func (Cmd) root() error {
	return os.Chdir("../")
}
