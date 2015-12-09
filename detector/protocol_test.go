// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMetric(t *testing.T) {
	line := "foo 1449655769 3.14"
	m, err := ParseMetric(line)
	assert.Nil(t, err)
	assert.Equal(t, m.Name, "foo")
	assert.Equal(t, m.Stamp, uint64(1449655769))
	assert.Equal(t, m.Value, 3.14)
}
