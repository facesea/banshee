// Copyright 2015 Eleme Inc. All rights reserved.

package statedb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestGirdNo(t *testing.T) {
	// New DB just for testing
	db := &DB{numGrid: 288, gridLen: 300}
	period := db.numGrid * db.gridLen
	// Consistence
	p := &models.Metric{Stamp: 1450428723}
	q := &models.Metric{Stamp: p.Stamp + uint32(period)}
	assert.Ok(t, db.getGridNo(p) == db.getGridNo(p))
	assert.Ok(t, db.getGridNo(q) == db.getGridNo(p))
	// Correcty
	m := &models.Metric{Stamp: 1450429041}
	n := &models.Metric{Stamp: m.Stamp + uint32(db.gridLen)}
	i := db.getGridNo(m)
	j := db.getGridNo(n)
	assert.Ok(t, i+1 == j || (i == db.numGrid-1 && j == 0))
}

func TestEncodeKey(t *testing.T) {
	// New DB just for testing
	db := &DB{numGrid: 288, gridLen: 300}
	period := db.numGrid * db.gridLen
	// Test
	m := &models.Metric{Name: "foo", Stamp: 1450429041}
	gridNo := db.getGridNo(m)
	s := fmt.Sprintf("foo:%d", gridNo)
	assert.Ok(t, s == string(db.encodeKey(m)))
	m.Stamp += uint32(period)
	assert.Ok(t, s == string(db.encodeKey(m)))
}

func TestEncodeValue(t *testing.T) {
	// New DB just for testing
	db := &DB{numGrid: 288, gridLen: 300}
	// Test
	s := &models.State{Average: 891.232898, StdDev: 1.2, Count: 123}
	v := "891.2329:1.2:123"
	assert.Ok(t, v == string(db.encodeValue(s)))
}

func TestDecodeValue(t *testing.T) {
	// New DB just for testing
	db := &DB{numGrid: 288, gridLen: 300}
	// Test
	value := []byte("123.12:1.23333:19")
	s, err := db.decodeValue(value)
	assert.Ok(t, err == nil)
	assert.Ok(t, s.Average == 123.12)
	assert.Ok(t, s.StdDev == 1.23333)
	assert.Ok(t, s.Count == 19)
}

func TestValueEncoding(t *testing.T) {
	// New DB just for testing
	db := &DB{numGrid: 288, gridLen: 300}
	// Test
	s := &models.State{Average: 182.092, StdDev: 1.3, Count: 18}
	value := db.encodeValue(s)
	n, err := db.decodeValue(value)
	assert.Ok(t, err == nil)
	assert.Ok(t, n.Average == s.Average)
	assert.Ok(t, n.StdDev == s.StdDev)
	assert.Ok(t, n.Count == s.Count)
}
