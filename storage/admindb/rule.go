// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// SetRuleChan set the changed rule channel
func (db *DB) SetRuleChan(rc chan *models.Rule) {
	db.ruleChan = rc
}

// changedRule send msg to cache that rule changed, cache should clean the
// dirty cache of the rule
func (db *DB) changedRule(rule *models.Rule) {
	select {
	case db.ruleChan <- rule:
	default:
		log.Error("buffered changed rules channel is full, cache won't clean the dirty cache..")
	}
}

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

// HasRule returns true if the rule of this id is in db.
func (db *DB) HasRule(id int) bool {
	_, ok := db.getRule(id)
	return ok
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
	db.rules.Put(r.ID, r)
	// send changed rule to cache
	db.changedRule(rule.Copy())
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
	// Get rule by id.
	rule, ok := db.getRule(id)
	if !ok {
		return ErrRuleNotFound
	}
	// Cache
	db.changedRule(rule.Copy())
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
