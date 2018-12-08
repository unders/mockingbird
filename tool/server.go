// +build mage

// Tests the github.com/unders/mockingbird code.
package main

import (
	"fmt"
	"os"
	"time"

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

// Build builds the mockingbird server here $GOPATH/bin/mockingbird
func (Server) Build() error {
	mg.SerialDeps(
		Server.All,
		Server.build,
	)

	return nil
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
	flags := fmt.Sprintf("-ldflags=%s", ldflags())
	cmd := "github.com/unders/mockingbird/server/cmd/mockingbird"

	fmt.Printf("go build -o %s %s %s\n", mockingbird, flags, cmd)
	return sh.RunV("go", "build", "-o", mockingbird, flags, cmd)
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

func ldflags() string {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	hash := hash()
	tag := tag()
	if tag == "" {
		tag = "dev"
	}
	return fmt.Sprintf(`-X "main.timestamp=%s" `+
		`-X "main.commitHash=%s" `+
		`-X "main.gitTag=%s"`, timestamp, hash, tag)
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		fmt.Println("ingoring fatal error; we set tag=dev")
	}
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}
