// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"runtime"
	"testing"
)

// Assert for testing.
func Assert(t *testing.T, b bool) {
	if !b {
		_, fileName, line, _ := runtime.Caller(1)
		t.Errorf("\n => %s:%d", fileName, line)
	}
}
