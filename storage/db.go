// Copyright 2015 Eleme Inc. All rights reserved.

// Package storage implements persistence storage based on leveldb.
package storage

import (
	"fmt"
	"github.com/eleme/banshee/storage/admindb"
	"github.com/eleme/banshee/storage/metricdb"
	"github.com/eleme/banshee/storage/statedb"
	"os"
	"path"
)

// DB file mode.
const filemode = 0755

// Child db filename.
const (
	admindbFileName  = "admin"
	metricdbFileName = "metric"
	statedbFileName  = "state"
)

// Metrics expiration = metricExpirationNumPeriods * period
const metricExpirationNumPeriods = 7

// Options is for db opening.
type Options struct {
	// sdb
	NumGrid int
	GridLen int
}

// DB handles the storage on leveldb.
type DB struct {
	Admin  *admindb.DB
	Metric *metricdb.DB
	State  *statedb.DB
}

// Open a DB by fileName and options.
func Open(fileName string, options *Options) (*DB, error) {
	// Create if not exist
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		err := os.Mkdir(fileName, filemode)
		if err != nil {
			return nil, err
		}
	}
	// Open databases.
	db := new(DB)
	db.Admin, err = admindb.Open(path.Join(fileName, admindbFileName))
	if err != nil {
		return nil, err
	}
	db.Metric, err = metricdb.Open(path.Join(fileName, metricdbFileName))
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s-%dx%d", statedbFileName, options.NumGrid, options.GridLen)
	opts := &statedb.Options{NumGrid: options.NumGrid, GridLen: options.GridLen}
	db.State, err = statedb.Open(path.Join(fileName, name), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close a DB.
func (db *DB) Close() error {
	if err := db.Admin.Close(); err != nil {
		return err
	}
	if err := db.Metric.Close(); err != nil {
		return err
	}
	if err := db.State.Close(); err != nil {
		return err
	}
	return nil
}
