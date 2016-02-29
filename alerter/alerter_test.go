// Copyright 2016 Eleme Inc. All rights reserved.

package alerter

import (
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestHourInRange(t *testing.T) {
	assert.Ok(t, hourInRange(3, 0, 6))
	assert.Ok(t, !hourInRange(7, 0, 6))
	assert.Ok(t, !hourInRange(6, 0, 6))
	assert.Ok(t, hourInRange(23, 19, 10))
	assert.Ok(t, hourInRange(6, 19, 10))
	assert.Ok(t, !hourInRange(13, 19, 10))
}
