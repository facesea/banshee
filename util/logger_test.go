// Copyright 2015 Eleme Inc. All rights reserved.

package util

import "testing"

func TestLogger(t *testing.T) {
	l := NewLogger("example")
	l.Debug("this is a debug message %v", 1)
	l.Info("this is a debug message %v", 2)
	l.Warn("this is a debug message %v", 3)
	l.Error("this is a debug message %v", 4)
}
