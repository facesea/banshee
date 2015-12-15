// Copyright 2015 Eleme Inc. All rights reserved.

// Package storage implements persistence storage based on leveldb.
package storage

import (
	"fmt"
	"github.com/eleme/banshee/storage/adb"
	"github.com/eleme/banshee/storage/mdb"
	"github.com/eleme/banshee/storage/sdb"
	"os"
	"path"
)

// DB file mode.
const filemode = 0755

// Child db filename.
const (
	adbFileName = "admin"
	mdbFileName = "metrics"
	sdbFileName = "states"
)

// DB opening options.
type Options struct {
	// sdb
	NumGrid int
	GridLen int
}

// DB handles the storage on leveldb.
type DB struct {
	adb *adb.DB
	mdb *mdb.DB
	sdb *sdb.DB
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
	db.adb, err = adb.Open(path.Join(fileName, adbFileName))
	if err != nil {
		return nil, err
	}
	db.mdb, err = mdb.Open(path.Join(fileName, mdbFileName))
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s-%dx%d", sdbFileName, options.NumGrid, options.GridLen)
	opts := &sdb.Options{NumGrid: options.NumGrid, GridLen: options.GridLen}
	db.sdb, err = sdb.Open(path.Join(fileName, name), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close a DB.
func (db *DB) Close() error {
	if err := db.adb.Close(); err != nil {
		return err
	}
	if err := db.mdb.Close(); err != nil {
		return err
	}
	if err := db.sdb.Close(); err != nil {
		return err
	}
	return nil
}

// Using adb handle.
func (db *DB) UsingA() *adb.DB {
	return db.adb
}

// Using mdb handle.
func (db *DB) UsingM() *mdb.DB {
	return db.mdb
}

// Using sdb handle.
func (db *DB) UsingS() *sdb.DB {
	return db.sdb
}
