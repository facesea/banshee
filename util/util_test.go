// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"github.com/eleme/banshee/util/assert"
	"os"
	"testing"
)

func TestMatch(t *testing.T) {
	assert.Ok(t, Match("abcdefg", "a*cd*fg"))
	assert.Ok(t, !Match("cbcdefg", "a*cd*fg"))
	assert.Ok(t, !Match("abcdef", "a*cd*fg"))
	assert.Ok(t, !Match("abxdef", "a*cd*fg"))
}

func TestIsFileExist(t *testing.T) {
	assert.Ok(t, IsFileExist(os.Args[0]))
	assert.Ok(t, !IsFileExist("file-not-exist"))
}
