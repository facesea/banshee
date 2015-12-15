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
	level int       = INFO
	w     io.Writer = os.Stderr
)

// Set logging name.
func SetName(s string) {
	name = s
}

// Set logging level.
func SetLevel(l int) {
	level = l % len(levelNames)
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

func Debug(format string, a ...interface{}) {
	log(DEBUG, format, a...)
}

func Info(format string, a ...interface{}) {
	log(INFO, format, a...)
}

func Warn(format string, a ...interface{}) {
	log(WARN, format, a...)
}

func Error(format string, a ...interface{}) {
	log(ERROR, format, a...)
}

func Fatal(format string, a ...interface{}) {
	log(ERROR, format, a...)
	os.Exit(1)
}
