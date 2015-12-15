// Copyright 2015 Eleme Inc. All rights reserved.

// Package sdb handles the states storage.
package sdb

import (
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
)

// DB opening options.
type Options struct {
	NumGrid int
	GridLen int
}

// DB handles the states storage for detection.
type DB struct {
	// LevelDB
	db *leveldb.DB
	// Period
	numGrid int
	gridLen int
}

// Open a DB by fileName and options.
func Open(fileName string, options *Options) (*DB, error) {
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		db:      db,
		numGrid: options.NumGrid,
		gridLen: options.GridLen,
	}, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// Operations

// Get the State for the given Metric.
func (db *DB) Get(m *models.Metric) (*models.State, error) {
	key := db.encodeKey(m)
	value, err := db.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return db.decodeValue(value)
}

// Put the State for the given Metric.
func (db *DB) Put(m *models.Metric, s *models.State) error {
	key := db.encodeKey(m)
	value := db.encodeValue(s)
	return db.db.Put(key, value, nil)
}
