// Copyright 2015 Eleme Inc. All rights reserved.

package cache

import "github.com/eleme/banshee/models"

// getRule returns rule by id.
func (c *Cache) getRule(id int) (*models.Rule, bool) {
	v, ok := c.rules.Get(id)
	if !ok {
		return nil, false
	}
	rule := v.(*models.Rule)
	return rule, true
}

// NumRules returns the number of rules.
func (c *Cache) NumRules() int {
	return c.rules.Len()
}

// GetRule returns rule.
func (c *Cache) GetRule(rule *models.Rule) error {
	r, ok := c.getRule(rule.ID)
	if !ok {
		return ErrRuleNotFound
	}
	r.CopyTo(rule)
	return nil
}

// HasRule checks whether a rule exist.
func (c *Cache) HasRule(rule *models.Rule) bool {
	return c.rules.Has(rule.ID)
}

// Rules returns all rules.
func (c *Cache) Rules(rules *[]*models.Rule) {
	for _, v := range c.rules.Items() {
		rule := v.(*models.Rule)
		*rules = append(*rules, rule.Copy())
	}
}

// RulesN returns rules for given range.
func (c *Cache) RulesN(rules *[]*models.Rule, offset int, limit int) {
	for _, v := range c.rules.ItemsN(offset, limit) {
		rule := v.(*models.Rule)
		*rules = append(*rules, rule.Copy())
	}
}

// AddRule adds a rule to cache.
func (c *Cache) AddRule(rule *models.Rule) {
	r := rule.Copy()
	r.MakeShared()
	c.rules.Put(rule.ID, r)
}

// DeleteRule deletes a rule from cache.
func (c *Cache) DeleteRule(rule *models.Rule) error {
	// Check.
	_, ok := c.getRule(rule.ID)
	if !ok {
		return ErrRuleNotFound
	}
	// Delete rule from its project.
	p, ok := c.getProject(rule.ID)
	if !ok {
		return ErrProjectNotFound
	}
	p.DeleteRule(rule.ID)
	// Delete r.
	if !c.rules.Delete(rule.ID) {
		return ErrRuleNotFound
	}
	return nil
}
