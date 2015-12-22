// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// getRule returns the rule by id.
func (db *DB) getRule(id int) (*models.Rule, bool) {
	v, ok := db.rules.Get(id)
	if !ok {
		// Not found
		return nil, false
	}
	rule := v.(*models.Rule)
	return rule, true
}

// Rules returns all the rules.
func (db *DB) Rules() (l []*models.Rule) {
	for _, v := range db.rules.Items() {
		rule := v.(*models.Rule)
		l = append(l, rule.Copy())
	}
	return l
}

// GetRule returns rule by id.
func (db *DB) GetRule(id int) (*models.Rule, error) {
	rule, ok := db.getRule(id)
	if !ok {
		return nil, ErrNotFound
	}
	return rule.Copy(), nil
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
	// Add to its project.
	proj, ok := db.getProject(rule.ProjectID)
	if !ok {
		return ErrNotFound
	}
	proj.AddRule(rule)
	// Marke as shared.
	rule.MakeShared()
	// Add to rules.
	db.rules.Set(rule.ID, rule)
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
	rule, ok := db.getRule(id)
	if !ok {
		return ErrNotFound
	}
	// Delete rule from its projects.
	rule = rule.Copy()
	proj, ok := db.getProject(rule.ProjectID)
	if !ok {
		return ErrNotFound
	}
	proj.DeleteRule(id)
	// Delete from rules.
	if !db.rules.Delete(id) {
		return ErrNotFound
	}
	return nil
}
