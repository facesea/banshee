// Copyright 2015 Eleme Inc. All rights reserved.

// Tool to migrate admin rules, projects and users from bell.
//
// Requirements
//
//	bell.js v2.0+ https://github.com/eleme/bell.js
//	banshee v0.0.7+ https://github.com/eleme/banshee
//
// Command Line Usage
//
//	./migrate -from bell.db -to banshee.db -with-users
//	mv banshee.db path/to/storage/admin
//
// Warning: Must backup your db files before the migration.
//
package main

import (
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
	withUsers         = flag.Bool("with-users", false, "if migrate users")
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
	// Bell DB
	if db, err := gorm.Open("sqlite3", *bellDBFileName); err != nil {
		log.Fatal("%s: %v", *bellDBFileName, err)
	} else {
		bellDB = &db
		bellDB.LogMode(false)
	}
	// Banshee DB
	if db, err := gorm.Open("sqlite3", *bansheeDBFileName); err != nil {
		log.Fatal("%s:%v", *bansheeDBFileName, err)
	} else {
		bansheeDB = &db
		bansheeDB.LogMode(false)
	}
	if err := bansheeDB.AutoMigrate(&models.Project{}, &models.Rule{}, &models.User{}).Error; err != nil {
		log.Fatal("failed to migrate schema for %s: %v", *bansheeDBFileName, err)
	}
}

// Main
//
//	1. Migrate projects and all their rules.
//	2. Migrate users and their user-project relationships.
//
func main() {
	migrateProjects()
	if *withUsers {
		migrateUsers()
	}
}

// Migrate projects.
//
//	1. Fetch all projects from belldb.
//	2. Create the project into bansheedb.
//	3. Create the rules for each project.
//
func migrateProjects() {
	var projs []Project
	// Fetch all projects from belldb.
	if err := bellDB.Find(&projs).Error; err != nil {
		log.Fatal("fetch all projects from %s: %v", *bellDBFileName, err)
	}
	for _, proj := range projs {
		// Create banshee project.
		if err := models.ValidateProjectName(proj.Name); err != nil {
			log.Warn("project %s: %v, skipping..", proj.Name, err)
			continue
		}
		p := &models.Project{Name: proj.Name}
		if err := bansheeDB.Create(p).Error; err != nil {
			sqliteErr, ok := err.(sqlite3.Error)
			if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				log.Warn("project %s already in %s, skipping..", p.Name, *bansheeDBFileName)
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
			// Create banshee rule.
			if err := models.ValidateRulePattern(rule.Pattern); err != nil {
				log.Warn("rule %s: %v, belongs to %s, skippig..", rule.Pattern, err, proj.Name)
				continue
			}
			r := &models.Rule{
				Pattern:   rule.Pattern,
				ProjectID: p.ID,
				TrendUp:   rule.Up,
				TrendDown: rule.Down,
				// Important: max and min for bell is reversed with banshee's.
				ThresholdMax: rule.Min,
				ThresholdMin: rule.Max,
			}
			if err := bansheeDB.Create(r).Error; err != nil {
				sqliteErr, ok := err.(sqlite3.Error)
				if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					log.Warn("rule %s already in %s, skipping..", r.Pattern, *bansheeDBFileName)
				} else {
					log.Fatal("cannot create rule %s: %v", r.Pattern, err)
				}
			}
		}
	}
}

// Migrate users.
//
//	1. Fetch all users from belldb.
//	2. Create the users into bansheedb.
//	3. Establish the relationships between project and user.
//
func migrateUsers() {
	var users []Receiver
	// Fetch all users from belldb.
	if err := bellDB.Find(&users).Error; err != nil {
		log.Fatal("fetch all users from %s: %v", *bellDBFileName, err)
	}
	for _, user := range users {
		// Create banshee user.
		err := models.ValidateUserName(user.Name)
		if err == nil {
			err = models.ValidateUserEmail(user.Email)
		}
		if err == nil {
			err = models.ValidateUserPhone(user.Phone)
		}
		if err != nil {
			log.Warn("user %s: %v, skipping..", user.Name, err)
		}
		u := &models.User{
			Name:        user.Name,
			Email:       user.Email,
			Phone:       user.Phone,
			EnableEmail: user.EnableEmail,
			EnablePhone: user.EnablePhone,
			Universal:   user.Universal,
		}
		if err := bansheeDB.Create(u).Error; err != nil {
			sqliteErr, ok := err.(sqlite3.Error)
			if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				log.Warn("user %s already in %s, skipping..", u.Name, *bansheeDBFileName)
			} else {
				log.Fatal("cannot create user %s: %v", u.Name, err)
			}
		}
		// Establish relationship to project.
		if user.Universal {
			continue
		}
		// Get all relationships for this user.
		var relations []ReceiverProject
		if err := bellDB.Where("ReceiverId = ?", user.ID).Find(&relations).Error; err != nil {
			log.Fatal("cannot fetch user-project relations for user %s: %v", user.Name, err)
		}
		for _, relation := range relations {
			var proj Project
			if err := bellDB.First(&proj, relation.ProjectID).Error; err != nil {
				if err == gorm.RecordNotFound {
					log.Warn("project %d not found for user %s, skipping..", relation.ProjectID, user.Name)
					continue
				}
				log.Fatal("cannot get project %d for user %s", relation.ProjectID, user.Name)
			}
			p := &models.Project{}
			if err := bansheeDB.Where("name = ?", proj.Name).First(p).Error; err != nil {
				if err == gorm.RecordNotFound {
					log.Warn("project %s not found in %s, skipping..", proj.Name, *bansheeDBFileName)
					continue
				}
				log.Fatal("cannot get project %s in %s", proj.Name, *bansheeDBFileName)
			}
			if err := bansheeDB.Model(p).Association("Users").Append(u).Error; err != nil {
				if err == gorm.RecordNotFound {
					log.Warn("record not found: %v", err)
					continue
				}
				log.Fatal("cannot append user %s to project %s:%v", u.Name, p.Name, err)
			}
		}
	}
}
