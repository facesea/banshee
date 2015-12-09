// Copyright 2015 Eleme Inc. All rights reserved.
package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCorrupted(t *testing.T) {
	assert.False(t, IsErrCorrupted(errors.New("something wrong")))
	assert.True(t, IsErrCorrupted(NewErrCorruptedWithString("something wrong")))
}
