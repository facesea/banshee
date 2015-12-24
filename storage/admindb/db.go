// Copyright 2015 Eleme Inc. All rights reserved.

// Package admindb handles the admin storage.
package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage/admindb/cache"
	"github.com/eleme/banshee/storage/admindb/persist"
)

// DB handles admin storage.
type DB struct {
	// Cache
	cache *cache.Cache
	// Persist
	persist *persist.Persist
	// Signals
	ruleChanges []chan *models.Rule
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
	if err := db.cache.Init(db.persist); err != nil {
		return nil, err
	}
	return db, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.persist.Close()
}
