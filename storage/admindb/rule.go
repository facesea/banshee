// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
)

// ListenRuleChanges adds a channel to listen rule changes.
func (db *DB) ListenRuleChanges(ch chan *models.Rule) {
	db.ruleChanges = append(db.ruleChanges, ch)
}

// ruleChanged send msg to cache that rule changed, cache should clean the
// dirty cache of the rule
func (db *DB) ruleChanged(rule *models.Rule) {
	for _, ch := range db.ruleChanges {
		select {
		case ch <- rule:
		default:
			log.Error("buffered changed rules channel is full, skipping..")
		}
	}
}

// NumRules returns the number of rules.
func (db *DB) NumRules() int {
	return db.cache.NumRules()
}

// Rules returns all rules.
func (db *DB) Rules(rules *[]*models.Rule) {
	db.cache.Rules(rules)
}

// RulesN returns rules for given range.
func (db *DB) RulesN(rules *[]*models.Rule, offset int, limit int) {
	db.cache.RulesN(rules, offset, limit)
}

// HasRule returns true if the rule of this id is in db.
func (db *DB) HasRule(rule *models.Rule) bool {
	return db.cache.HasRule(rule)
}

// GetRule returns a rule.
func (db *DB) GetRule(rule *models.Rule) error {
	return db.cache.GetRule(rule)
}

// AddRuleToProject adds a rule to cache and adds to project.
func (db *DB) AddRuleToProject(proj *models.Project, rule *models.Rule) error {
	rule.ProjectID = proj.ID
	if err := db.persist.AddRule(rule); err != nil {
		return err
	}
	db.cache.AddRule(rule)
	if err := db.cache.AddRuleToProject(proj, rule); err != nil {
		return err
	}
	db.ruleChanged(rule)
	return nil
}

// DeleteRule deletes a rule from the db.
func (db *DB) DeleteRule(rule *models.Rule) error {
	if err := db.cache.GetRule(rule); err != nil {
		return err
	}
	if err := db.persist.DeleteRule(rule); err != nil {
		return err
	}
	if err := db.cache.DeleteRule(rule); err != nil {
		return err
	}
	db.ruleChanged(rule)
	return nil
}
