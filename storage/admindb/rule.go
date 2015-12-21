// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// Rules returns all the rules.
func (db *DB) Rules() (l []*models.Rule) {
	for _, rule := range db.rules.Items() {
		l = append(l, rule.(*models.Rule))
	}
	return l
}

// GetRule returns rule by id.
func (db *DB) GetRule(id int) (*models.Rule, error) {
	v, ok := db.rules.Get(id)
	if !ok {
		return nil, ErrNotFound
	}
	return v.(*models.Rule), nil
}

// AddRule adds a rule to the db.
func (db *DB) AddRule(rule *models.Rule) error {
	// Sql
	if err := db.db.Create(rule).Error; err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		return err
	}
	// Cache
	// Add to rules.
	db.rules.Set(rule.ID, rule)
	// Add to its project.
	proj, err := db.GetProject(rule.ProjectID)
	if err != nil {
		return err
	}
	proj.AddRule(rule.Clone())
	return nil
}

// DeleteRule deletes a rule from the db.
func (db *DB) DeleteRule(id int) error {
	// Sql.
	if err := db.db.Delete(&models.Rule{ID: id}).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		return err
	}
	// Cache
	// Get rule by id.
	rule, err := db.GetRule(id)
	if err != nil {
		return err
	}
	// Delete from its projects.
	proj, err := db.GetProject(rule.ProjectID)
	if err != nil {
		return err
	}
	proj.DeleteRule(id)
	// Delete from rules.
	if !db.rules.Delete(id) {
		return ErrNotFound
	}
	return nil
}
