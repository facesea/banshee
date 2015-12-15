// Copyright 2015 Eleme Inc. All rights reserved.

// Package assert provides a method to assert booleans.
package assert

import (
	"runtime"
	"testing"
)

// Assert the given boolean is True.
func Ok(t *testing.T, b bool) {
	if !b {
		_, fileName, line, _ := runtime.Caller(1)
		t.Errorf("\nassertion failed:%s:%d", fileName, line)
	}
}
