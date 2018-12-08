package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/domain/mockingbird/mock"
	"github.com/unders/mockingbird/server/pkg/signal"
)

//
// From the build script
//
var (
	timestamp  = "timestamp"
	commitHash = "commitHash"
	gitTag     = "gitTag"
)

func options() Options {
	var addr = ":8080"
	flag.StringVar(&addr, "http.addr", addr, "HTTP address.")
	flag.Parse()

	l := log.New(os.Stderr, "", 0)
	return Options{
		Timestamp:  timestamp,
		CommitHash: commitHash,
		GitTag:     gitTag,

		ServerAddr:              addr,
		ServerReadHeaderTimeout: 30 * time.Second,  // 30s
		ServerReadTimeout:       50 * time.Second,  // 50s
		ServerWriteTimeout:      240 * time.Second, // 240s = 4*60s => 4 minutes
		ServerIdleTimeout:       20 * time.Second,  // 20
		ServerShutdownTimeout:   300 * time.Second, // 300s = 5*60s => 5 minutes

		StartTime: time.Now().UTC(),
		Log:       &mockingbird.Logger{Log: l},
		ErrorLog:  l,
	}
}

func main() {
	if err := run(options()); err != nil {
		os.Exit(1)
	}
}

func run(o Options) error {
	l := o.Log

	format := "mockingbird server is starting commit-hash=%s build-time=%s start-time=%s "
	l.Info(fmt.Sprintf(format, o.CommitHash, o.Timestamp, o.StartTime.Format(time.RFC3339)))
	format = "Options%+v"
	l.Info(fmt.Sprintf(format, o))

	handler := &ochttp.Handler{
		Handler: Handler{
			HTML: mock.HTMLAdapter{Code: 200, Body: []byte("Hello World")},
			Log:  o.Log,
		}.Make(),

		Propagation: &b3.HTTPFormat{},
	}

	srv := &http.Server{
		Handler:           handler,
		Addr:              o.ServerAddr,
		ReadHeaderTimeout: o.ServerReadHeaderTimeout,
		ReadTimeout:       o.ServerReadTimeout,
		WriteTimeout:      o.ServerWriteTimeout,
		IdleTimeout:       o.ServerIdleTimeout,
		ErrorLog:          o.ErrorLog,
	}

	format = "mockingbird listens on addr %s commit-hash=%s run-time=%s"
	l.Info(fmt.Sprintf(format, o.ServerAddr, o.CommitHash, time.Since(o.StartTime)))

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	var err error
	select {
	case err = <-errCh:
		const format = "server error=%s commit-hash=%s run-time=%s"
		l.Error(fmt.Sprintf(format, err, o.CommitHash, time.Since(o.StartTime)))
	case sig := <-signal.Interrupt():
		const format = "got interrupt signal=%s commit-hash=%s run-time=%s"
		l.Info(fmt.Sprintf(format, sig, o.CommitHash, time.Since(o.StartTime)))
	}

	stopTime := time.Now().UTC()
	waitTimeout := o.ServerShutdownTimeout
	format = "shutting down the http server commit-hash=%s wait-timeout=%s run-time=%s"
	l.Info(fmt.Sprintf(format, o.CommitHash, waitTimeout, time.Since(o.StartTime)))
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		const format = "server shutdown error=%s commit-hash=%s shutdown-time=%s"
		l.Error(fmt.Sprintf(format, err, o.CommitHash, time.Since(stopTime)))
	}

	format = "mockingbird server is stopped commit-hash=%s shutdown-time=%s run-time=%s"
	l.Info(fmt.Sprintf(format, o.CommitHash, time.Since(stopTime), time.Since(o.StartTime)))

	return err
}
