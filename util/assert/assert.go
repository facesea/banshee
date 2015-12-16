// Copyright 2015 Eleme Inc. All rights reserved.

// Package assert provides a method to assert booleans.
package assert

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Ok asserts that the given boolean is True.
func Ok(t *testing.T, b bool) {
	if !b {
		_, fileName, line, _ := runtime.Caller(1)
		cwd, err := os.Getwd()
		if err != nil {
			t.Errorf("unexcepted:%v", err)
		}
		fileName, err = filepath.Rel(cwd, fileName)
		if err != nil {
			t.Errorf("unexcepted:%v", err)
		}
		t.Errorf("\nassertion failed: %s:%d", fileName, line)
	}
}
