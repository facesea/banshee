// Copyright 2015 Eleme Inc. All rights reserved.

package statedb

import (
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Options is for db opening.
type Options struct {
	NumGrid uint32
	GridLen uint32
}

// DB handles the states storage for detection.
type DB struct {
	// LevelDB
	db *leveldb.DB
	// Period
	numGrid uint32
	gridLen uint32
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

// Delete all states for a metric name.
// This operation is currently only used for cleaning.
func (db *DB) Delete(name string) error {
	// Name must be the key prefix
	iter := db.db.NewIterator(util.BytesPrefix([]byte(name)), nil)
	batch := new(leveldb.Batch)
	for iter.Next() {
		key := iter.Key()
		batch.Delete(key)
	}
	if batch.Len() > 0 {
		return db.db.Write(batch, nil)
	}
	return nil
}
