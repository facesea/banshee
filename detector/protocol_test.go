// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/util"
	"testing"
)

func TestParseMetric(t *testing.T) {
	line := "foo 1449655769 3.14"
	m, err := parseMetric(line)
	util.Assert(t, err == nil)
	util.Assert(t, m.Name == "foo")
	util.Assert(t, m.Stamp == uint32(1449655769))
	util.Assert(t, m.Value == 3.14)
}
