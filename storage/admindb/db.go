// Copyright 2015 Eleme Inc. All rights reserved.

// Package admindb handles the admin storage.
package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage/admindb/rcache"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Import but no use
)

const dialect = "sqlite3"

// DB handles admin storage.
type DB struct {
	// DB
	db *gorm.DB
	// Cache
	rc *rcache.RCache
}

// Open DB by fileName.
func Open(fileName string) (*DB, error) {
	// Open
	gdb, err := gorm.Open(dialect, fileName)
	if err != nil {
		return nil, err
	}
	db := new(DB)
	db.db = &gdb
	// Migration
	if err := db.migrate(); err != nil {
		return nil, err
	}
	// Cache
	db.rc = rcache.New()
	if err := db.rc.Init(db.db); err != nil {
		return nil, err
	}
	return db, nil
}

// Close DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// DB returns db handle.
func (db *DB) DB() *gorm.DB {
	return db.db
}

// migrate db schema.
func (db *DB) migrate() error {
	rule := &models.Rule{}
	user := &models.User{}
	proj := &models.Project{}
	return db.db.AutoMigrate(rule, user, proj).Error
}

// RulesCache handle.
func (db *DB) RulesCache() *rcache.RCache {
	return db.rc
}
