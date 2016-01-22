// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

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
	db, err := Open(fileName)
	assert.Ok(t, err == nil)
	assert.Ok(t, util.IsFileExist(fileName))
	db.Close()
	os.RemoveAll(fileName)
}

func TestPut(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Put.
	m := &models.Metric{
		Name:    "foo",
		Stamp:   1452758773,
		Value:   3.14,
		Score:   0.1892,
		Average: 3.133,
	}
	err := db.Put(m)
	assert.Ok(t, err == nil)
	// Must in db
	key := encodeKey(m)
	value, err := db.db.Get(key, nil)
	assert.Ok(t, err == nil)
	m1 := &models.Metric{
		Name:  m.Name,
		Stamp: m.Stamp,
	}
	err = decodeValue(value, m1)
	assert.Ok(t, err == nil)
	assert.Ok(t, reflect.DeepEqual(m, m1))
}

func TestGet(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Nothing.
	ms, err := db.Get("foo", 0, 1452758773)
	assert.Ok(t, err == nil)
	assert.Ok(t, len(ms) == 0)
	// Put some.
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758723})
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758733, Value: 1.89, Score: 1.12, Average: 1.72})
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758743})
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758753})
	// Get again.
	ms, err = db.Get("foo", 1452758733, 1452758753)
	assert.Ok(t, err == nil)
	assert.Ok(t, len(ms) == 2)
	// Test the value.
	m := ms[0]
	assert.Ok(t, m.Value == 1.89 && m.Score == 1.12)
}

func TestDelete(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Nothing.
	err := db.Delete("foo", 0, 1452758773)
	assert.Ok(t, err == nil)
	// Put some.
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758723})
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758733})
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758743})
	db.Put(&models.Metric{Name: "foo", Stamp: 1452758753})
	// Delete again
	err = db.Delete("foo", 1452758733, 1452758753)
	assert.Ok(t, err == nil)
	// Get
	ms, err := db.Get("foo", 1452758723, 1452758763)
	assert.Ok(t, len(ms) == 2)
	assert.Ok(t, ms[0].Stamp == 1452758723)
	assert.Ok(t, ms[1].Stamp == 1452758753)
}
