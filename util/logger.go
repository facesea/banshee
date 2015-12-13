// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// Levels
const (
	LOG_DEBUG int = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
)

// Level names
var logLevelNames = [4]string{"DEBUG", "INFO", "WARN", "ERROR"}

// Logger is a logging handler.
type Logger struct {
	name  string
	level int
	w     io.Writer
}

// NewLogger creates a new Logger, the default logging level is LOG_INGO and
// default writer is os.Stderr.
func NewLogger(name string) *Logger {
	l := new(Logger)
	l.name = name
	l.SetLevel(LOG_INFO)
	l.SetWriter(os.Stderr)
	return l
}

// SetName sets the name for a logger.
func (l *Logger) SetName(name string) {
	l.name = name
}

// SetLevel sets the level for a logger.
func (l *Logger) SetLevel(level int) {
	l.level = level % len(logLevelNames)
}

// SetWriter sets the writer for a logger.
func (l *Logger) SetWriter(w io.Writer) {
	l.w = w
}

func (l *Logger) doLog(level int, format string, a ...interface{}) {
	if level >= l.level {
		msg := fmt.Sprintf(format, a...)
		now := time.Now().String()[:23]
		pid := os.Getpid()
		s := fmt.Sprintf("%s %s %s[%d]: %s", now, logLevelNames[level], l.name, pid, msg)
		fmt.Fprintln(l.w, s)
	}
}

// Log a message with debug level.
func (l *Logger) Debug(format string, a ...interface{}) {
	l.doLog(LOG_DEBUG, format, a...)
}

// Log a message with info level.
func (l *Logger) Info(format string, a ...interface{}) {
	l.doLog(LOG_INFO, format, a...)
}

// Log a message with warn level.
func (l *Logger) Warn(format string, a ...interface{}) {
	l.doLog(LOG_WARN, format, a...)
}

// Log a message with error level.
func (l *Logger) Error(format string, a ...interface{}) {
	l.doLog(LOG_ERROR, format, a...)
}

// Log a message with error level and terminate the process.
func (l *Logger) Fatal(format string, a ...interface{}) {
	l.doLog(LOG_ERROR, format, a...)
	os.Exit(1)
}

// Log current running env
func (l *Logger) Runtime() {
	l.doLog(LOG_INFO, "runtime info: using %s, up to %d cpus..", runtime.Version(), runtime.GOMAXPROCS(0))
}

// The global logger
var logger = NewLogger("util")

// Get the root logger.
func GetRootLogger() *Logger {
	return logger
}
