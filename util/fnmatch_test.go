// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"testing"
)

func TestFnMatch(t *testing.T) {
	Assert(t, FnMatch("abcdefg", "a*cd*fg"))
	Assert(t, !FnMatch("cbcdefg", "a*cd*fg"))
	Assert(t, !FnMatch("abcdef", "a*cd*fg"))
	Assert(t, !FnMatch("abxdef", "a*cd*fg"))
}
