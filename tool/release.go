package main

import (
	"errors"
	"os"

	"github.com/magefile/mage/sh"
)

// Release creates a release: TAG=v1.2.3 mage release
func Release() (err error) {
	if os.Getenv("TAG") == "" {
		return errors.New("TAG environment variable is required")
	}
	if err = sh.RunV("git", "tag", "-a", "$TAG"); err != nil {
		return err
	}
	if err = sh.RunV("git", "push", "origin", "$TAG"); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = sh.RunV("git", "tag", "--delete", "$TAG")
			_ = sh.RunV("git", "push", "--delete", "origin", "$TAG")
		}
	}()

	return sh.RunV("goreleaser")
}
