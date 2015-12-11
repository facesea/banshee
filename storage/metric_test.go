// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"os"
	"reflect"
	"testing"
)

func TestMetric(t *testing.T) {
	cfg := config.NewWithDefaults()
	cfg.Storage.Path = "storage_test/"
	db, err := Open(cfg)
	util.Assert(t, err == nil)
	defer db.Close()
	defer os.RemoveAll(cfg.Storage.Path)
	m := &models.Metric{Name: "foo", Stamp: 1449740827, Value: 1.43}
	err = db.GetMetric(m)
	util.Assert(t, err == ErrNotFound)
	n := &models.Metric{Name: "foo", Stamp: 1449740827, Value: 1.43, Avg: 1.4, Std: 0.1, Count: 2}
	err = db.PutMetric(n)
	util.Assert(t, err == nil)
	err = db.GetMetric(m)
	util.Assert(t, err == nil)
	util.Assert(t, reflect.DeepEqual(n, m))
}
