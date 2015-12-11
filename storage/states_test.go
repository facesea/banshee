// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"os"
	"reflect"
	"testing"
)

func TestState(t *testing.T) {
	fileName := "storage_test/"
	numGrids, gridLen := 480, 180
	db, err := Open(fileName, numGrids, gridLen)
	util.Assert(t, err == nil)
	defer db.Close()
	defer os.RemoveAll(fileName)
	m := &models.Metric{Name: "foo", Stamp: 1449740827, Value: 1.43}
	s, err := db.GetState(m)
	util.Assert(t, err == ErrNotFound)
	s = &models.State{Average: 1.4, StdDev: 0.1, Count: 2}
	err = db.PutState(m, s)
	util.Assert(t, err == nil)
	s1, err := db.GetState(m)
	util.Assert(t, err == nil)
	util.Assert(t, reflect.DeepEqual(s, s1) && s != s1)
}
