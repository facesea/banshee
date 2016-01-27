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
	"flag"
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
	bansheeDB.AutoMigrate(&models.Project{}, &models.Rule{}, &models.User{})
}

func main() {
	// Fixme
}
