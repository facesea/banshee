// Copyright 2015 Eleme Inc. All rights reserved.

// Package persist hanldles admin persistence storage.
package persist

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
)

const dialect = "sqlite3"

// Persist
type Persist struct {
	// db handle.
	db gorm.DB
}

// Open a DB by fileName.
func Open(fileName string) (*Persist, error) {
	// Open db.
	db, err := gorm.Open(dialect, fileName)
	if err != nil {
		return nil, err
	}
	p := new(Persist)
	p.db = db
	// Migrate schema.
	err = persist.Migrate()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Close the db.
func (p *Persist) Close() error {
	return p.db.Close()
}

// Migrate checks db schema and do auto migratation.
func (p *Persist) Migrate() error {
	rule := &models.Rule{}
	user := &models.User{}
	proj := &models.Project{}
	return p.db.AutoMigrate(rule, user, proj).Error
}
