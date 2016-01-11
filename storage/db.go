// Copyright 2015 Eleme Inc. All rights reserved.

// Package storage implements persistence storage based on leveldb.
package storage

import (
	"fmt"
	"github.com/eleme/banshee/storage/admindb"
	"github.com/eleme/banshee/storage/indexdb"
	"github.com/eleme/banshee/storage/metricdb"
	"github.com/eleme/banshee/storage/statedb"
	"github.com/eleme/banshee/util/log"
	"os"
	"path"
)

// DB file mode.
const filemode = 0755

// Child db filename.
const (
	admindbFileName  = "admin"
	indexdbFileName  = "index"
	metricdbFileName = "metric"
	statedbFileName  = "state"
)

// Options is for db opening.
type Options struct {
	// statedb
	NumGrid uint32
	GridLen uint32
}

// DB handles the storage on leveldb.
type DB struct {
	// Child db
	Admin  *admindb.DB
	Index  *indexdb.DB
	Metric *metricdb.DB
	State  *statedb.DB
}

// Open a DB by fileName and options.
func Open(fileName string, options *Options) (*DB, error) {
	// Create if not exist
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Debug("create dir %s", fileName)
		err := os.Mkdir(fileName, filemode)
		if err != nil {
			return nil, err
		}
	}
	// Admindb.
	db := new(DB)
	db.Admin, err = admindb.Open(path.Join(fileName, admindbFileName))
	if err != nil {
		return nil, err
	}
	// Indexdb.
	db.Index, err = indexdb.Open(path.Join(fileName, indexdbFileName))
	if err != nil {
		return nil, err
	}
	// Metricdb.
	db.Metric, err = metricdb.Open(path.Join(fileName, metricdbFileName))
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s-%dx%d", statedbFileName, options.NumGrid, options.GridLen)
	opts := &statedb.Options{NumGrid: options.NumGrid, GridLen: options.GridLen}
	// Statedb.
	db.State, err = statedb.Open(path.Join(fileName, name), opts)
	if err != nil {
		return nil, err
	}
	log.Debug("storage is opened successfully")
	return db, nil
}

// Close a DB.
func (db *DB) Close() error {
	// Admindb.
	if err := db.Admin.Close(); err != nil {
		return err
	}
	// Indexdb.
	if err := db.Index.Close(); err != nil {
		return err
	}
	// Metricdb.
	if err := db.Metric.Close(); err != nil {
		return err
	}
	// Statedb.
	if err := db.State.Close(); err != nil {
		return err
	}
	return nil
}
