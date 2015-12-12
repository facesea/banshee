// Copyright 2015 Eleme Inc. All rights reserved.

package util

import "testing"

func TestLogger(t *testing.T) {
	l := NewLogger("example")
	l.Runtime()
	l.Debug("this is a debug message %v", 1)
	l.Info("this is a info message %v", 2)
	l.Warn("this is a warn message %v", 3)
	l.Error("this is a error message %v", 4)
}
