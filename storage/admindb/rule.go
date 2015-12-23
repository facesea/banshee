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

// GetRule returns rule into a local value.
func (db *DB) GetRule(r *models.Rule) error {
	rule, ok := db.getRule(r.ID)
	if !ok {
		return ErrRuleNotFound
	}
	rule.CopyTo(r)
	return nil
}

// AddRuleToProject adds a rule to a project.
func (db *DB) AddRuleToProject(proj *models.Project, rule *models.Rule) error {
	// If proj exist.
	_, ok := db.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	// Create projectID
	rule.ProjectID = proj.ID
	// Sql: rule.ID will be created.
	if err := db.db.Create(rule).Error; err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		if err == sqlite3.ErrConstraintNotNull {
			return ErrConstraintNotNull
		}
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		return err
	}
	// Append rule to proj.Rules.
	proj.AddRule(rule)
	// Cache a copy.
	r := rule.Copy()
	// Add to its project.
	p, ok := db.getProject(r.ProjectID)
	if !ok {
		return ErrProjectNotFound
	}
	p.AddRule(r)
	// Mark as shared.
	r.MakeShared()
	// Add to rules.
	db.rules.Set(r.ID, r)
	return nil
}

// DeleteRule deletes a rule from the db.
func (db *DB) DeleteRule(id int) error {
	// Sql.
	if err := db.db.Delete(&models.Rule{ID: id}).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrRuleNotFound
		}
		return err
	}
	// Cache
	// Get rule by id.
	rule, ok := db.getRule(id)
	if !ok {
		return ErrRuleNotFound
	}
	// Delete rule from its projects.
	proj, ok := db.getProject(rule.GetProjectID())
	if !ok {
		return ErrProjectNotFound
	}
	proj.DeleteRule(id)
	// Delete from rules.
	if !db.rules.Delete(id) {
		return ErrRuleNotFound
	}
	return nil
}
