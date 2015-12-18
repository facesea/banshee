// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestEncodeKey(t *testing.T) {
	m := &models.Metric{Name: "foo", Stamp: horizon + 0xf}
	key := encodeKey(m)
	s := "foo000000f"
	assert.Ok(t, s == string(key))
}

func TestDecodeKey(t *testing.T) {
	key := []byte("foo000001f")
	m := &models.Metric{}
	err := decodeKey(key, m)
	assert.Ok(t, err == nil)
	assert.Ok(t, m.Name == "foo")
	assert.Ok(t, m.Stamp == 36+0xf+horizon)
}

func TestStampLenEnoughToUse(t *testing.T) {
	stamp := uint32(90*365*24*60*60) + horizon
	m := &models.Metric{Name: "foo", Stamp: stamp}
	key := encodeKey(m)
	n := &models.Metric{}
	err := decodeKey(key, n)
	assert.Ok(t, err == nil)
	assert.Ok(t, n.Name == m.Name)
	assert.Ok(t, n.Stamp == m.Stamp)
}

func TestEncodeValue(t *testing.T) {
	m := &models.Metric{Value: 1.23, Score: 0.72, Average: 0.798766}
	value := encodeValue(m)
	s := "1.23:0.72:0.79877"
	assert.Ok(t, s == string(value))
}

func TestDecodeValue(t *testing.T) {
	m := &models.Metric{}
	value := []byte("1.23:0.72:0.79")
	err := decodeValue(value, m)
	assert.Ok(t, err == nil)
	assert.Ok(t, m.Value == 1.23)
	assert.Ok(t, m.Score == 0.72)
	assert.Ok(t, m.Average == 0.79)
}
