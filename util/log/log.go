// Copyright 2015 Eleme Inc. All rights reserved.

// Package log implements leveled logging.
package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

// Level
const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
)

// Level name
var levelNames = [4]string{"DEBUG", "INFO", "WARN", "ERROR"}

// Logging runtime
var (
	name  string
	level           = INFO
	w     io.Writer = os.Stderr
)

// SetName sets the logging name.
func SetName(s string) {
	name = s
}

// SetLevel sets the logging level.
func SetLevel(l int) {
	level = l % len(levelNames)
}

// SetWriter sets the writer.
func SetWriter(writer io.Writer) {
	w = writer
}

// Do logging.
func log(l int, format string, a ...interface{}) {
	if l >= level {
		// Caller location
		_, fileName, line, _ := runtime.Caller(2)
		dir := path.Dir(fileName)
		base := path.Base(fileName)
		loc := path.Join(path.Base(dir), base)
		// Meta
		msg := fmt.Sprintf(format, a...)
		now := time.Now().String()[:23]
		pid := os.Getpid()
		s := fmt.Sprintf("%s %-5s %s[%d] <%s:%d>: %s", now, levelNames[l], name, pid, loc, line, msg)
		fmt.Fprintln(w, s)
	}
}

// Debug logs message with level DEBUG.
func Debug(format string, a ...interface{}) {
	log(DEBUG, format, a...)
}

// Info logs message with level INFO.
func Info(format string, a ...interface{}) {
	log(INFO, format, a...)
}

// Warn logs message with level WARN.
func Warn(format string, a ...interface{}) {
	log(WARN, format, a...)
}

// Error logs message with level ERROR.
func Error(format string, a ...interface{}) {
	log(ERROR, format, a...)
}

// Fatal logs message with level FATAL.
func Fatal(format string, a ...interface{}) {
	log(ERROR, format, a...)
	os.Exit(1)
}
