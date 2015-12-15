// Copyright 2015 Eleme Inc. All rights reserved.

// Package mdb handles the storage for metrics.
package mdb

import "github.com/syndtr/goleveldb/leveldb"

// DB handles the metric storage.
type DB struct {
	// LevelDB
	db *leveldb.DB
}

// Open a DB by fileName.
func Open(fileName string) (*DB, error) {
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}
