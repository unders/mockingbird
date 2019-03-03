// +build mage

// Tests the github.com/unders/mockingbird code.
package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Server mg.Namespace

// All runs all server test targets: server:test and server:race
func (Server) All() {
	mg.SerialDeps(
		Go.Version,
		Server.download,
		Server.Test,
		Server.Race,
	)
}

// Test runs go test github.com/unders/mockingbird/cmd/...
func (Server) Test() error {
	mg.SerialDeps(Server.download)
	return sh.RunV("go", "test", "./server/...")
}

// Race runs go test github.com/unders/mockingbird/cmd/...
func (Server) Race() error {
	mg.SerialDeps(Server.download)
	return sh.RunV("go", "test", "-race", "./server/...")
}

// Build builds the mockingbird server into dir $GOPATH/bin/mockingbird
func (Server) Build() error {
	mg.SerialDeps(
		Server.All,
		Server.build,
	)

	return nil
}

// Start starts $GOPATH/bin/mockingbird
func (Server) Start() error {
	mg.SerialDeps(
		Server.All,
		Server.build,
	)

	return sh.RunV("mockingbird", "-l")
}

func (Server) build() error {
	gopath := ""
	if home := os.Getenv("HOME"); home != "" {
		gopath = fmt.Sprintf("%s/go", home)
	}
	if path := os.Getenv("GOPATH"); path != "" {
		gopath = path
	}
	mockingbird := fmt.Sprintf("%s/bin/mockingbird", gopath)
	cmd := "github.com/unders/mockingbird/server/cmd/mockingbird"

	fmt.Printf("go build -o %s %s \n", mockingbird, cmd)
	return sh.RunV("go", "build", "-o", mockingbird, cmd)
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
