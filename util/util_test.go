// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"github.com/eleme/banshee/util/assert"
	"os"
	"testing"
)

func TestToFixed(t *testing.T) {
	assert.Ok(t, ToFixed(1.2345, 2) == "1.23")
	assert.Ok(t, ToFixed(10000.12121121, 5) == "10000.12121")
}

func TestIsFileExist(t *testing.T) {
	assert.Ok(t, IsFileExist(os.Args[0]))
	assert.Ok(t, !IsFileExist("file-not-exist"))
}
