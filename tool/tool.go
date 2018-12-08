// +build mage

package main

import (
	"github.com/magefile/mage/sh"
	"github.com/unders/mockingbird/tool/pkg/env"
)

// Tools ensures we have a sync'd all our tools
func Tools() error {
	// mg.Deps(checkProtoc)

	update, err := env.Bool("Update")
	if err != nil {
		return err
	}

	retool := "github.com/twitchtv/retool"
	args := []string{"get", retool}
	if update {
		args = []string{"get", "-u", retool}
	}

	if err := sh.Run("go", args...); err != nil {
		return err
	}

	return sh.Run("retool", "sync")
}
