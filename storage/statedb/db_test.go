// Copyright 2015 Eleme Inc. All rights reserved.

package statedb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/assert"
	"os"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	fileName := "db-testing"
	db, err := Open(fileName, &Options{288, 300})
	assert.Ok(t, err == nil)
	assert.Ok(t, util.IsFileExist(fileName))
	db.Close()
	os.RemoveAll(fileName)
}

func TestPut(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName, &Options{288, 300})
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Test.
	m := &models.Metric{Name: "foo", Stamp: 1450435291}
	s := &models.State{Count: 1, Average: 100, StdDev: 1.02}
	err := db.Put(m, s)
	assert.Ok(t, err == nil)
	// Must in db
	key := db.encodeKey(m)
	value, err := db.db.Get(key, nil)
	assert.Ok(t, err == nil)
	s1, _ := db.decodeValue(value)
	assert.Ok(t, reflect.DeepEqual(s1, s))
}

func TestGet(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName, &Options{288, 300})
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Not found.
	m := &models.Metric{Name: "foo", Stamp: 1450435291}
	_, err := db.Get(m)
	assert.Ok(t, err == ErrNotFound)
	// Put one.
	s := &models.State{Count: 1, Average: 100, StdDev: 1.02}
	db.Put(m, s)
	// Get again, must in db.
	s1, err := db.Get(m)
	assert.Ok(t, err == nil)
	assert.Ok(t, reflect.DeepEqual(s1, s))
}

func TestDelete(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName, &Options{288, 300})
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Put some states.
	m := &models.Metric{Name: "bar", Stamp: 1450435742}
	n := &models.Metric{Name: "foo", Stamp: 1450435742}
	db.Put(m, &models.State{})
	db.Put(n, &models.State{})
	db.Put(&models.Metric{Name: m.Name, Stamp: 1450435741}, &models.State{})
	// Get one.
	s, err := db.Get(m)
	assert.Ok(t, err == nil && s != nil)
	// Delete one.
	err = db.Delete(m.Name)
	assert.Ok(t, err == nil)
	// Get agian
	_, err = db.Get(m)
	assert.Ok(t, err == ErrNotFound)
	s1, err := db.Get(n)
	assert.Ok(t, err == nil && s1 != nil)
}
