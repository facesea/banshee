// Copyright 2015 Eleme Inc. All rights reserved.

// Package admindb handles the admin storage.
package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Import but no use
)

const dialect = "sqlite3"

// DB handles admin storage.
type DB struct {
	// DB
	db *gorm.DB
	// Cache
	RulesCache *rulesCache
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
	db.RulesCache = newRulesCache()
	if err := db.RulesCache.Init(db.db); err != nil {
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
	log.Debug("migrate db schema..")
	rule := &models.Rule{}
	user := &models.User{}
	proj := &models.Project{}
	return db.db.AutoMigrate(rule, user, proj).Error
}
