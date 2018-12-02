package signal

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/unders/mockingbird/cmd/pkg/testdata"
)

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	t.Helper()
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("\nWant: %v\n Got: %v\n", sig, s)
		}
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("\nWant: %s\n Got: timeout error\n", sig)
	}
}

func TestInterrupt(t *testing.T) {
	c := Interrupt()
	defer signal.Stop(c)

	// Send this process a SIGTERM
	t.Log("sigterm...")
	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	testdata.AssertNil(t, err)
	waitSig(t, c, syscall.SIGTERM)

	// Send this process a SIGINT
	t.Log("sigint...")
	err = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	testdata.AssertNil(t, err)
	waitSig(t, c, syscall.SIGINT)
}
