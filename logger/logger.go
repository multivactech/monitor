package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jrick/logrotate/rotator"
	"github.com/multivactech/monitor/btclog"
)

var (
	logRotator *rotator.Rotator
	Log        btclog.Logger
)

type logWriter struct{}

func (logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	if logRotator != nil {
		logRotator.Write(p)
	}
	return len(p), nil
}

// logCleanup does the necessary cleaning before system shuts down.
func LogCleanup() {
	if logRotator != nil {
		logRotator.Close()
	}
}

// InitLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func InitLogRotator(logFile string) {
	logDir, _ := filepath.Split(logFile)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		os.Exit(1)
	}
	os.Remove(logFile)
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		os.Exit(1)
	}

	logRotator = r
}

func logInit() {
	backendLog := btclog.NewBackend(logWriter{})
	Log = backendLog.Logger("")
}

func SetLogInfo() {
	logInit()
	Log.SetLevel(btclog.LevelInfo)
}

func SetLogDebug() {
	logInit()
	Log.SetLevel(btclog.LevelDebug)
}

func SetLogError() {
	logInit()
	Log.SetLevel(btclog.LevelError)
}

func SetLogWarn() {
	logInit()
	Log.SetLevel(btclog.LevelWarn)
}
