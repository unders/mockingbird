// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	"github.com/unders/mockingbird/tool/pkg/env"
)

type Tool mg.Namespace

// Download download toll required by this project
func (Tool) Download() error {
	mg.Deps(
	// install dependencies in parallel
	// checkProtoc
	// goreleaser
	// gometalinter
	)

	return errors.New("not implemented")

	update, err := env.Bool("Update")
	if err != nil {
		return err
	}

	//
	// Example of how to install a go tool dependecy
	//
	retool := "github.com/twitchtv/retool"
	args := []string{"get", retool}
	if update {
		args = []string{"get", "-u", retool}
	}
	return sh.RunV("go", args...)
}
