// Copyright 2015 Eleme Inc. All rights reserved.

// Package admindb handles the admin storage.
package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/safemap"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Import but not use
)

// DB handles admin storage.
type DB struct {
	// Gorm
	db gorm.DB
	// Cache
	projects *safemap.SafeMap
	rules    *safemap.SafeMap
	users    *safemap.SafeMap
}

// Open a DB by fileName.
func Open(fileName string) (*DB, error) {
	sdb, err := gorm.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	db := new(DB)
	db.db = sdb
	db.projects = safemap.New()
	db.rules = safemap.New()
	db.users = safemap.New()
	if err := db.autoMigrate(); err != nil {
		return nil, err
	}
	return db, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// autoMigrate creates tables if there are not exist.
func (db *DB) autoMigrate() error {
	return db.db.AutoMigrate(&models.Project{}, &models.User{}, &models.Rule{}).Error
}
