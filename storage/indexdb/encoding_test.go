// Copyright 2015 Eleme Inc. All rights reserved.

package indexdb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	idx := &models.Index{Stamp: 1450422576, Score: 1.2, Average: 3.1}
	s := "1450422576:1.20000:3.10000"
	assert.Ok(t, string(encode(idx)) == s)
}

func TestDecode(t *testing.T) {
	s := "1450422576:1.20000:3.10000"
	idx := &models.Index{}
	err := decode([]byte(s), idx)
	assert.Ok(t, err == nil)
	assert.Ok(t, idx.Stamp == 1450422576)
	assert.Ok(t, idx.Score == 1.2)
	assert.Ok(t, idx.Average == 3.1)
}
