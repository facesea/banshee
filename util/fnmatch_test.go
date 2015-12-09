// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFnMatch(t *testing.T) {
	assert.True(t, FnMatch("abcdefg", "a*cd*fg"))
	assert.False(t, FnMatch("cbcdefg", "a*cd*fg"))
	assert.False(t, FnMatch("abcdef", "a*cd*fg"))
	assert.False(t, FnMatch("abxdef", "a*cd*fg"))
}
