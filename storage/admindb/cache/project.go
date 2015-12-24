// Copyright 2015 Eleme Inc. All rights reserved.

package cache

import "github.com/eleme/banshee/models"

// getProject returns project by id.
func (c *Cache) getProject(id int) (*models.Project, bool) {
	v, ok := c.projs.Get(id)
	if !ok {
		return nil, false
	}
	proj := v.(*models.Project)
	return proj, true
}

// NumProjects returns the number of projects.
func (c *Cache) NumProjects() int {
	return c.projs.Len()
}

// GetProject returns project.
func (c *Cache) GetProject(proj *models.Project) error {
	p, ok := c.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	p.CopyTo(proj)
	return nil
}

// HasProject checks whether a project exist.
func (c *Cache) HasProject(proj *models.Project) bool {
	return c.projs.Has(proj.ID)
}

// GetProjects returns all projects.
func (c *Cache) GetProjects(projs *[]*models.Project) {
	for _, v := range c.projs.Items() {
		proj := v.(*models.Project)
		*projs = append(*projs, proj.Copy())
	}
}

// GetProjectsN returns projects for given range.
func (c *Cache) GetProjectsN(projs *[]*models.Project, offset int, limit int) {
	for _, v := range c.projs.ItemsN(offset, limit) {
		proj := v.(*models.Project)
		*projs = append(*projs, proj.Copy())
	}
}

// AddProject adds a project to cache.
func (c *Cache) AddProject(proj *models.Project) {
	copy := proj.Copy()
	copy.MakeShared()
	c.projs.Put(proj.ID, copy)
}

// UpdateProject updates a project.
func (c *Cache) UpdateProject(proj *models.Project) error {
	// Find
	p, ok := c.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	// Update
	p.Update(proj)
	// Update users
	users := p.GetUsers()
	for _, user := range users {
		u, ok := c.getUser(user.ID)
		if !ok {
			return ErrUserNotFound
		}
		if !u.UpdateProject(p) {
			return ErrProjectNotFound
		}
	}
	return nil
}

// DeleteProject deletes a project.
func (c *Cache) DeleteProject(proj *models.Project) error {
	// Find.
	p, ok := c.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	// Delete its rules.
	rules := p.GetRules()
	for _, rule := range rules {
		if !c.rules.Delete(rule.ID) {
			return ErrRuleNotFound
		}
	}
	// Delete projects from its users.
	users := p.GetUsers()
	for _, user := range users {
		u, ok := c.getUser(user.ID)
		if !ok {
			return ErrUserNotFound
		}
		if u.DeleteProject(proj.ID) {
			return ErrProjectNotFound
		}
	}
	// Delete p
	if !c.projs.Delete(proj.ID) {
		return ErrProjectNotFound
	}
	return nil
}

// AddRuleToProject adds a rule to a project.
func (c *Cache) AddRuleToProject(proj *models.Project, rule *models.Rule) error {
	// Check proj.
	p, ok := c.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	// Check rule.
	r, ok := c.getRule(rule.ID)
	if !ok {
		return ErrRuleNotFound
	}
	// Add r to p
	p.AddRule(r)
	// Add rule to proj.
	proj.AddRule(rule)
	return nil
}

// AddUserToProject adds a user to a project.
func (c *Cache) AddUserToProject(proj *models.Project, user *models.User) error {
	// Check proj.
	p, ok := c.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	// Check user.
	u, ok := c.getUser(user.ID)
	if !ok {
		return ErrUserNotFound
	}
	// Add u to p.
	p.AddUser(u)
	// Add p to u.
	u.AddProject(p)
	// Add proj to user.
	if !user.HasProject(proj.ID) {
		user.AddProject(proj)
	}
	// Add user to proj.
	if !proj.HasUser(user.ID) {
		proj.AddUser(user)
	}
	return nil
}
