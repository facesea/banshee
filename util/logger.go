// Copyright 2015 Eleme Inc. All rights reserved.

// Custom logger to stderr
package util

import (
	"fmt"
	"io"
	"os"
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

type Logger struct {
	name  string
	level int
	w     io.Writer
}

func NewLogger(name string) *Logger {
	l := new(Logger)
	l.name = name
	l.SetLevel(LOG_INFO)
	l.SetWriter(os.Stderr)
	return l
}

func (l *Logger) SetLevel(level int) {
	l.level = level % 4
}

func (l *Logger) SetWriter(w io.Writer) {
	l.w = w
}

func (l *Logger) doLog(level int, format string, a ...interface{}) {
	if level >= l.level {
		msg := fmt.Sprintf(format, a...)
		now := time.Now().String()[:23]
		pid := os.Getpid()
		s := fmt.Sprintf("%s %s %s[%d]: %s", now, logLevelNames[l.level], l.name, pid, msg)
		fmt.Fprintln(os.Stderr, s)
	}
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.doLog(LOG_DEBUG, format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.doLog(LOG_INFO, format, a...)
}

func (l *Logger) Warn(format string, a ...interface{}) {
	l.doLog(LOG_WARN, format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.doLog(LOG_ERROR, format, a...)
}

func (l *Logger) Fatal(format string, a ...interface{}) {
	l.doLog(LOG_ERROR, format, a...)
	os.Exit(1)
}
