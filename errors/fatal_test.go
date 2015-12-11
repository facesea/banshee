// Copyright 2015 Eleme Inc. All rights reserved.
package errors

import (
	"errors"
	"github.com/eleme/banshee/util"
	"testing"
)

func TestIsFatal(t *testing.T) {
	util.Assert(t, !IsFatal(errors.New("something wrong")))
	util.Assert(t, IsFatal(NewErrFatalWithString("something wrong")))
}
