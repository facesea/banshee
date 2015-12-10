// Copyright 2015 Eleme Inc. All rights reserved.
package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFatal(t *testing.T) {
	assert.False(t, IsFatal(errors.New("something wrong")))
	assert.True(t, IsFatal(NewErrFatalWithString("something wrong")))
}
