// Copyright 2015 Eleme Inc. All rights reserved.

// Package cache handles admin cache.
package cache

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage/admindb/persist"
	"github.com/eleme/banshee/util/skiplist"
)

// Cache is db cache.
type Cache struct {
	projs *skiplist.Skiplist
	rules *skiplist.Skiplist
	users *skiplist.Skiplist
}

// New creates cache.
func New() *Cache {
	c := new(Cache)
	c.projs = skiplist.New()
	c.rules = skiplist.New()
	c.users = skiplist.New()
	return c
}

// Clear the cache.
func (c *Cache) Clear() {
	c.projs.Clear()
	c.users.Clear()
	c.rules.Clear()
}

// Init cache from persist.
func (c *Cache) Init(p *persist.Persist) error {
	if err := c.initRules(p); err != nil {
		return err
	}
	if err := c.initUsers(p); err != nil {
		return err
	}
	if err := c.initProjects(p); err != nil {
		return err
	}
	return nil
}

// initRules inits rules from persist.
func (c *Cache) initRules(p *persist.Persist) error {
	var rules []*models.Rule
	if err := p.Rules(&rules); err != nil {
		return err
	}
	for _, rule := range rules {
		rule.MakeShared()
		c.rules.Put(rule.ID, rule)
	}
	return nil
}

// initUsers inits users from persist.
func (c *Cache) initUsers(p *persist.Persist) error {
	var users []*models.User
	if err := p.Users(&users); err != nil {
		return err
	}
	for _, user := range users {
		var projs []*models.Project
		if err := p.ProjectsOfUser(user, &projs); err != nil {
			return err
		}
		for _, proj := range projs {
			user.AddProject(proj)
		}
		user.MakeShared()
		c.users.Put(user.ID, user)
	}
	return nil
}

// initProjects inits projs from persist.
func (c *Cache) initProjects(p *persist.Persist) error {
	var projs []*models.Project
	if err := p.Projects(&projs); err != nil {
		return err
	}
	for _, proj := range projs {
		var rules []*models.Rule
		if err := p.RulesOfProject(proj, &rules); err != nil {
			return err
		}
		for _, rule := range rules {
			proj.AddRule(rule)
		}
		var users []*models.User
		if err := p.UsersOfProject(proj, &users); err != nil {
			return err
		}
		for _, user := range users {
			proj.AddUser(user)
		}
		proj.MakeShared()
		c.projs.Put(proj.ID, proj)
	}
	return nil
}
