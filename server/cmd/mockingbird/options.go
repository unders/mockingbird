package main

import (
	"log"
	"time"

	"github.com/unders/mockingbird/server/domain/mockingbird"
)

// Options defined required fields for mockingbird
type Options struct {
	// From build script
	Timestamp  string
	CommitHash string
	GitTag     string

	StartTime time.Time

	// http.Server settings
	ServerAddr              string
	ServerReadHeaderTimeout time.Duration
	ServerReadTimeout       time.Duration
	ServerWriteTimeout      time.Duration
	ServerIdleTimeout       time.Duration
	ServerShutdownTimeout   time.Duration

	Log      mockingbird.Log
	ErrorLog *log.Logger
}
