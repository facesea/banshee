// Copyright 2015 Eleme Inc. All rights reserved.

// Tool to migrate admin rules, projects and users from bell.
//
// Requirements
//
//	bell.js v2.0+ https://github.com/eleme/bell.js
//	banshee v0.07+ https://github.com/eleme/banshee
//
// Command Line Usage
//
//	./migrate -from bell.db -to banshee.db
//	mv banshee.db path/to/storage/admin
//
// Warning: Must backup your db files before the migration.
//
package main

import (
	// "errors"
	"flag"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

var (
	// Arguments
	bellDBFileName    = flag.String("from", "bell.db", "bell db file name")
	bansheeDBFileName = flag.String("to", "banshee.db", "banshee db file name")
	// DB Handles
	bellDB    *gorm.DB
	bansheeDB *gorm.DB
)

// Init.
//
//	1. Parse arguments.
//	2. Init db handles.
//	3. Auto migrate schema.
//
func init() {
	flag.Parse()
	log.Info("migrate from %s to %s..", *bellDBFileName, *bansheeDBFileName)
	// Bell DB
	if db, err := gorm.Open("sqlite3", *bellDBFileName); err != nil {
		log.Fatal("%s: %v", *bellDBFileName, err)
	} else {
		bellDB = &db
	}
	patchReceiverProjectsFieldNames(bellDB)
	// Banshee DB
	if db, err := gorm.Open("sqlite3", *bansheeDBFileName); err != nil {
		log.Fatal("%s:%v", *bansheeDBFileName, err)
	} else {
		bansheeDB = &db
	}
	if err := bansheeDB.AutoMigrate(&models.Project{}, &models.Rule{}, &models.User{}).Error; err != nil {
		log.Fatal("failed to migrate schema for %s: %v", *bansheeDBFileName, err)
	}
}

// Main
//
//	1. Migrate projects and all their rules.
//	2. Migrate users and establish relations to their projects.
//	3. Log failure rows to console.
//
func main() {
	migrateProjs()
}

// Migrate projects.
//
//	1. Fetch all projects from belldb.
//	2. Create the project into bansheedb.
//	3. Create the rules for each project.
//
func migrateProjs() {
	// Fetch all projects from belldb.
	var projs []Project
	if err := bellDB.Find(&projs).Error; err != nil {
		log.Fatal("failed to fetch all projects from %s: %v", *bellDBFileName, err)
	}
	for _, proj := range projs {
		// Create banshee project.
		p := &models.Project{Name: proj.Name}
		if err := bansheeDB.Create(p).Error; err != nil {
			sqliteErr, ok := err.(sqlite3.Error)
			if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				log.Warn("project %s already in %s", p.Name, *bansheeDBFileName)
			} else {
				log.Fatal("cannot create project %s: %v", p.Name, err)
			}
		}
		// Fetch its rules from belldb.
		var rules []Rule
		if err := bellDB.Model(proj).Related(&rules).Error; err != nil {
			log.Fatal("cannot fetch rules for %s: %v", p.Name, err)
		}
		for _, rule := range rules {
			// FIXME: Validate rule pattern.
			// Create banshee rule.
			r := &models.Rule{
				Pattern:   rule.Pattern,
				ProjectID: p.ID,
				TrendUp:   rule.Up,
				TrendDown: rule.Down,
				// Important: reverse min/max here -_#
				ThresholdMax: rule.Min,
				ThresholdMin: rule.Max,
			}
			if err := bansheeDB.Create(r).Error; err != nil {
				sqliteErr, ok := err.(sqlite3.Error)
				if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					log.Warn("rule %s already in %s", r.Pattern, *bansheeDBFileName)
				} else {
					log.Fatal("cannot create rule %s: %v", r.Pattern, err)
				}
			}
		}
	}
}
