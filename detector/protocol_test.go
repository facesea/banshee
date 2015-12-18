// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestParseMetric(t *testing.T) {
	line := "foo 1449655769 3.14"
	m, err := parseMetric(line)
	assert.Ok(t, err == nil)
	assert.Ok(t, m.Name == "foo")
	assert.Ok(t, m.Stamp == uint32(1449655769))
	assert.Ok(t, m.Value == 3.14)
}

func TestParseMetricBadLine(t *testing.T) {
	line := "foo 1.3 1.234"
	m, err := parseMetric(line)
	assert.Ok(t, err != nil)
	assert.Ok(t, m == nil)
}
