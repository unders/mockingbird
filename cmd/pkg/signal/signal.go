// Package signal listens for these OS signals:
//
//
// SIGTERM
//                 The SIGTERM signal is sent to a process to request
//                 its termination. Unlike the SIGKILL signal, it can
//                 be caught and interpreted or ignored by the process.
//                 This allows the process to perform nice termination
//                 releasing resources and saving state if appropriate.
//                 SIGINT is nearly identical to SIGTERM.
//
// SIGINT
//                 The SIGINT signal is sent to a process by its controlling
//                 terminal when a user wishes to interrupt the process. This
//                 is typically initiated by pressing Ctrl-C, but on some
//                 systems, the "delete" character or "break" key can be used.
//
package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Interrupt returns an os.Signal chan that waits for these signals:
//
//      * SIGTERM   - request for termination (process can release its resources).
//      * SIGINT    - sent from terminal when it wishes to interrupt the process.
//
// Usage:
//
//       sig := <-signal.Interrupt()
//
func Interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt,
		syscall.SIGTERM,
		syscall.SIGINT)

	return interrupt
}
