// Copyright 2015 Eleme Inc. All rights reserved.

package indexdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	idx := &models.Index{Stamp: 1450422576, Score: 1.2, Average: 3.1}
	s := "1450422576:1.2:3.1"
	fmt.Printf("\n%s", string(encode(idx)))
	assert.Ok(t, string(encode(idx)) == s)
}

func TestDecode(t *testing.T) {
	s := "1450422576:1.2:3.1"
	idx := &models.Index{}
	err := decode([]byte(s), idx)
	assert.Ok(t, err == nil)
	assert.Ok(t, idx.Stamp == 1450422576)
	assert.Ok(t, idx.Score == 1.2)
	assert.Ok(t, idx.Average == 3.1)
}

func TestEncoding(t *testing.T) {
	idx := &models.Index{Stamp: 1450426828, Score: 0.678888, Average: 877.234}
	value := encode(idx)
	idx1 := &models.Index{}
	err := decode(value, idx1)
	assert.Ok(t, err == nil)
	assert.Ok(t, idx1.Stamp == idx.Stamp)
	assert.Ok(t, idx1.Score == 0.67889)
	assert.Ok(t, idx1.Average == 877.234)
}
