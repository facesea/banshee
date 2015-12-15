// Copyright 2015 Eleme Inc. All rights reserved.

// Package db handles the administration storage.
package adb

import (
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
)

// DB handles the administration storage including rules, users etc.
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

// Operations

// Get all rules from memory.
func (db *DB) GetRules() []models.Rule {
	// FIXME
	return []models.Rule{
		models.Rule{
			Pattern: "*",
			When:    models.WhenTrendUp | models.WhenTrendDown,
		},
	}
}
