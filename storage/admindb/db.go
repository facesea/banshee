// Copyright 2015 Eleme Inc. All rights reserved.

// Package admindb handles the admin storage.
package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/safemap"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // Import but no use
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
	if err := db.load(); err != nil {
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

// load cache from db on bootstrap.
func (db *DB) load() error {
	// Rules
	if err := db.loadRules(); err != nil {
		return err
	}
	// Users
	if err := db.loadUsers(); err != nil {
		return err
	}
	// Projects
	if err := db.loadProjects(); err != nil {
		return err
	}
	return nil
}

// loadRules loads rules from db on bootstrap.
func (db *DB) loadRules() error {
	var rules []models.Rule
	if err := db.db.Find(&rules).Error; err != nil {
		return err
	}
	for _, rule := range rules {
		r := &rule
		// Add to cache.
		r.MakeShared()
		db.rules.Set(rule.ID, r)
	}
	return nil
}

// loadUsers loads users from db on bootstrap.
func (db *DB) loadUsers() error {
	var users []models.User
	if err := db.db.Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		u := &user
		// Projects
		var projs []models.Project
		if err := db.db.Model(u).Related(&projs, "Projects").Error; err != nil {
			return err
		}
		for _, proj := range projs {
			u.AddProject(&proj)
		}
		// Add to cache.
		u.MakeShared()
		db.users.Set(u.ID, u)
	}
	return nil
}

// loadProjects loads projects from db on bootstrap.
func (db *DB) loadProjects() error {
	var projs []models.Project
	if err := db.db.Find(&projs).Error; err != nil {
		return err
	}
	for _, proj := range projs {
		p := &proj
		// Rules
		var rules []models.Rule
		if err := db.db.Model(p).Related(&rules).Error; err != nil {
			return err
		}
		for _, rule := range rules {
			p.AddRule(&rule)
		}
		// Users
		var users []models.User
		if err := db.db.Model(p).Related(&users, "Users").Error; err != nil {
			return err
		}
		for _, user := range users {
			p.AddUser(&user)
		}
		// Add to cache.
		p.MakeShared()
		db.projects.Set(p.ID, p)
	}
	return nil
}
