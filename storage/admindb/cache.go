// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
	"github.com/eleme/banshee/util/safemap"
	"github.com/jinzhu/gorm"
)

type rulesCache struct {
	// Cache
	rules *safemap.SafeMap
	// Listeners
	lnsAdd []chan *models.Rule
	lnsDel []chan *models.Rule
}

// newRulesCache creates a rulesCache.
func newRulesCache() *rulesCache {
	c := new(rulesCache)
	c.rules = safemap.New()
	c.lnsAdd = make([]chan *models.Rule, 0)
	c.lnsDel = make([]chan *models.Rule, 0)
	return c
}

// Init cache from db.
func (c *rulesCache) Init(db *gorm.DB) error {
	// Query
	var rules []models.Rule
	err := db.Find(&rules).Error
	if err != nil {
		return err
	}
	// Load
	for _, rule := range rules {
		// Share
		r := &rule
		r.Share()
		c.rules.Set(rule.ID, r)
	}
	return nil
}

// Len returns the number of rules in cache.
func (c *rulesCache) Len() int {
	return c.rules.Len()
}

// Get returns rule.
func (c *rulesCache) Get(rule *models.Rule) bool {
	r, ok := c.rules.Get(rule.ID)
	if !ok {
		return false
	}
	r.(*models.Rule).CopyTo(rule)
	return true
}

// Put a rule into cache.
func (c *rulesCache) Put(rule *models.Rule) bool {
	if c.rules.Has(rule.ID) {
		return false
	}
	r := rule.Copy()
	r.Share()
	c.rules.Set(rule.ID, r)
	c.pushAdded(rule)
	return true
}

// All returns all rules.
func (c *rulesCache) All(rules *[]*models.Rule) {
	for _, v := range c.rules.Items() {
		rule := v.(*models.Rule)
		*rules = append(*rules, rule.Copy())
	}
}

// Delete a rule from cache.
func (c *rulesCache) Delete(rule *models.Rule) bool {
	r, ok := c.rules.Pop(rule.ID)
	if ok {
		c.pushDeled(r.(*models.Rule).Copy())
		return true
	}
	return false
}

// OnAdd listens rules changes.
func (c *rulesCache) OnAdd(ch chan *models.Rule) {
	c.lnsAdd = append(c.lnsAdd, ch)
}

// OnDel listens rules changes.
func (c *rulesCache) OnDel(ch chan *models.Rule) {
	c.lnsDel = append(c.lnsDel, ch)
}

// pushAdded pushes changed rule to listeners.
func (c *rulesCache) pushAdded(rule *models.Rule) {
	for _, ch := range c.lnsAdd {
		select {
		case ch <- rule:
		default:
			log.Error("buffered rule added is full, skipping..")
		}
	}
}

// pushDeled pushes changed rule to listeners.
func (c *rulesCache) pushDeled(rule *models.Rule) {
	for _, ch := range c.lnsDel {
		select {
		case ch <- rule:
		default:
			log.Error("buffered rule delete is full, skipping..")
		}
	}
}
