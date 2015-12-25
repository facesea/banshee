// Copyright 2015 Eleme Inc. All rights reserved.

// Package rcache handles rules cache.
package rcache

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/log"
	"github.com/eleme/banshee/util/safemap"
	"github.com/jinzhu/gorm"
)

// RCache is rules cache.
type RCache struct {
	// Cache
	rules *safemap.SafeMap
	// Listeners
	lns []chan *models.Rule
}

// New creates a RCache.
func New() *RCache {
	c := new(RCache)
	c.rules = safemap.New()
	c.lns = make([]chan *models.Rule, 0)
	return c
}

// Init cache from db.
func (c *RCache) Init(db *gorm.DB) error {
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
func (c *RCache) Len() int {
	return c.rules.Len()
}

// Get returns rule.
func (c *RCache) Get(rule *models.Rule) bool {
	r, ok := c.rules.Get(rule.ID)
	if !ok {
		return false
	}
	r.(*models.Rule).CopyTo(rule)
	return true
}

// Put a rule into cache.
func (c *RCache) Put(rule *models.Rule) bool {
	if c.rules.Has(rule.ID) {
		return false
	}
	r := rule.Copy()
	r.Share()
	c.rules.Set(rule.ID, r)
	c.pushChanged(rule)
	return true
}

// All returns all rules.
func (c *RCache) All(rules *[]*models.Rule) {
	for _, v := range c.rules.Items() {
		rule := v.(*models.Rule)
		*rules = append(*rules, rule.Copy())
	}
}

// Delete a rule from cache.
func (c *RCache) Delete(rule *models.Rule) bool {
	r, ok := c.rules.Pop(rule.ID)
	if ok {
		c.pushChanged(r.(*models.Rule).Copy())
		return true
	}
	return false
}

// OnChange listens rules changes.
func (c *RCache) OnChange(ch chan *models.Rule) {
	c.lns = append(c.lns, ch)
}

// pushChanged pushes changed rule to listeners.
func (c *RCache) pushChanged(rule *models.Rule) {
	for _, ch := range c.lns {
		select {
		case ch <- rule:
		default:
			log.Error("buffered rule changes is full, skipping..")
		}
	}
}
