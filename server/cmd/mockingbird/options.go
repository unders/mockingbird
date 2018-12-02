package main

import (
	"log"
	"time"

	"github.com/unders/mockingbird/server/domain/mockinbird"
)

// Options defined required fields for mockingbird
type Options struct {
	Version    string
	Buildstamp string
	Githash    string

	StartTime time.Time

	// http.Server settings
	ServerAddr              string
	ServerReadHeaderTimeout time.Duration
	ServerReadTimeout       time.Duration
	ServerWriteTimeout      time.Duration
	ServerIdleTimeout       time.Duration
	ServerShutdownTimeout   time.Duration

	Log      mockinbird.Log
	ErrorLog *log.Logger
}
