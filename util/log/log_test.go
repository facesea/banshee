// Copyright 2015 Eleme Inc. All rights reserved.

package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	SetLevel(DEBUG)
	Debug("hello %s", "world")
	Info("hello %s", "world")
	Warn("hello %s", "world")
	Error("hello %s", "world")
}
