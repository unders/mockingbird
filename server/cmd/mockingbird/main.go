package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/unders/mockingbird/server/domain/mockingbird/app"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"

	"github.com/unders/mockingbird/server/domain/mockingbird"
	"github.com/unders/mockingbird/server/pkg/signal"
)

//
// From the build script
//

func options() Options {
	var (
		addr  = ":8080"
		local = false
	)
	flag.StringVar(&addr, "http.addr", addr, "HTTP address.")
	flag.BoolVar(&local, "l", local, "if app is running on a local dev server")
	flag.Parse()

	env := mockingbird.Env(os.Getenv("ENVIRONMENT"))
	if local {
		env = mockingbird.DEV
	}

	l := log.New(os.Stderr, "", 0)
	return Options{
		Env: env,

		ServerAddr:              addr,
		ServerReadHeaderTimeout: 30 * time.Second,  // 30s
		ServerReadTimeout:       50 * time.Second,  // 50s
		ServerWriteTimeout:      240 * time.Second, // 240s = 4*60s => 4 minutes
		ServerIdleTimeout:       20 * time.Second,  // 20
		ServerShutdownTimeout:   300 * time.Second, // 300s = 5*60s => 5 minutes

		FaviconDir:  "../web/mockingbird/favicon",
		TemplateDir: "../web/mockingbird/tmpl",
		AssetDir:    "../web/mockingbird/public",

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

	format := "mockingbird server is starting  start-time=%s "
	l.Info(fmt.Sprintf(format, o.StartTime.Format(time.RFC3339)))
	format = "Options%+v"
	l.Info(fmt.Sprintf(format, o))

	builder, err := app.Create(app.Options{
		Env:         o.Env,
		Logger:      o.ErrorLog,
		FaviconDir:  o.FaviconDir,
		TemplateDir: o.TemplateDir,
		AssetDir:    o.AssetDir,
	})
	if err != nil {
		return errors.Wrap(err, "app.Create() failed")
	}

	h := &ochttp.Handler{
		Handler: createHandler(handler{
			Favicon: builder.Favicon(),
			Assets:  builder.Assets(),
			HTML:    builder.HTMLAdapter(),
			Log:     builder.Log(),
		}),

		Propagation: &b3.HTTPFormat{},
	}

	srv := &http.Server{
		Handler:           h,
		Addr:              o.ServerAddr,
		ReadHeaderTimeout: o.ServerReadHeaderTimeout,
		ReadTimeout:       o.ServerReadTimeout,
		WriteTimeout:      o.ServerWriteTimeout,
		IdleTimeout:       o.ServerIdleTimeout,
		ErrorLog:          o.ErrorLog,
	}

	format = "mockingbird listens on addr %s run-time=%s"
	l.Info(fmt.Sprintf(format, o.ServerAddr, time.Since(o.StartTime)))

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case err = <-errCh:
		const format = "server error=%s run-time=%s"
		l.Error(fmt.Sprintf(format, err, time.Since(o.StartTime)))
	case sig := <-signal.Interrupt():
		const format = "got interrupt signal=%s run-time=%s"
		l.Info(fmt.Sprintf(format, sig, time.Since(o.StartTime)))
	}

	stopTime := time.Now().UTC()
	waitTimeout := o.ServerShutdownTimeout
	format = "shutting down the http server wait-timeout=%s run-time=%s"
	l.Info(fmt.Sprintf(format, waitTimeout, time.Since(o.StartTime)))
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		const format = "server shutdown error=%s shutdown-time=%s"
		l.Error(fmt.Sprintf(format, err, time.Since(stopTime)))
	}

	format = "mockingbird server is stopped shutdown-time=%s run-time=%s"
	l.Info(fmt.Sprintf(format, time.Since(stopTime), time.Since(o.StartTime)))

	return err
}
