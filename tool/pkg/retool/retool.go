package retool

import "github.com/magefile/mage/sh"

// ShellRun runs a command using a retool-cached binary.
func ShellRun(cmd string, args ...string) error {
	return sh.Run("retool", append([]string{"do", cmd}, args...)...)
}
