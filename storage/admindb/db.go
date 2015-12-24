// Copyright 2015 Eleme Inc. All rights reserved.

// Package admindb handles the admin storage.
package admindb

import (
	"github.com/eleme/banshee/storage/admindb/cache"
	"github.com/eleme/banshee/storage/admindb/persist"
)

// DB handles admin storage.
type DB struct {
	cache   *cache.Cache
	persist *persist.Persist
	// changed rule channel
	ruleChan chan *models.Rule
}

// Open a DB by fileName.
func Open(fileName string) (*DB, error) {
	db := new(DB)
	db.cache = cache.New()
	var err error
	db.persist, err = persist.Open(fileName)
	if err != nil {
		return nil, err
	}
	// Init cache.
	db.cache.Init(db.persist)
	return db, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.persist.Close()
}
