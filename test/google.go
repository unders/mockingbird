// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Google mg.Namespace

// Test tests the google service
func (Google) Test() error {
	return sh.RunV("go", "test", "-count=1", "github.com/unders/mockingbird/test/service/google/...")
}
