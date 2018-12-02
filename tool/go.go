// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Go mg.Namespace

// Version runs: go version
func (Go) Version() {
	sh.RunV("go", "version")
}
