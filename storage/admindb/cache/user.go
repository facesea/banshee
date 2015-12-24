// Copyright 2015 Eleme Inc. All rights reserved.

package cache

import "github.com/eleme/banshee/models"

// getUser returns user by id.
func (c *Cache) getUser(id int) (*models.User, bool) {
	v, ok := c.users.Get(id)
	if !ok {
		return nil, false
	}
	user := v.(*models.User)
	return user, true
}

// GetUser returns user.
func (c *Cache) GetUser(user *models.User) error {
	u, ok := c.getUser(user.ID)
	if !ok {
		return ErrUserNotFound
	}
	u.CopyTo(user)
	return nil
}

// HasUser checks whether a user exist.
func (c *Cache) HasUser(user *models.User) bool {
	return c.users.Has(user.ID)
}

// Users returns all users.
func (c *Cache) Users(users *[]*models.User) {
	for _, v := range c.users.Items() {
		user := v.(*models.User)
		*users = append(*users, user.Copy())
	}
}

// UsersN returns users for given range.
func (c *Cache) UsersN(users *[]*models.User, offset int, limit int) {
	for _, v := range c.users.ItemsN(offset, limit) {
		user := v.(*models.User)
		*users = append(*users, user.Copy())
	}
}

// AddUser adds a user to cache.
func (c *Cache) AddUser(user *models.User) {
	copy := user.Copy()
	copy.MakeShared()
	c.users.Put(user.ID, copy)
}

// UpdateUser updates a user.
func (c *Cache) UpdateUser(user *models.User) {
	// Find
	u, ok := c.getUser(user.ID)
	if !ok {
		return ErrUserNotFound
	}
	// Update
	u.Update(user)
	// Update projs
	projs := u.GetProjects()
	for _, proj := range projs {
		p, ok := c.getProject(proj.ID)
		if !ok {
			return ErrProjectNotFound
		}
		if !p.UpdateUser(u) {
			return ErrUserNotFound
		}
	}
	return nil
}

// DeleteUser deletes a user from cache.
func (c *Cache) DeleteUser(user *models.User) error {
	// Check
	u, ok := c.getUser(user.ID)
	if !ok {
		return ErrUserNotFound
	}
	// Delete from its projects.
	projs := user.GetProjects()
	for _, proj := range projs {
		p, ok := c.getProject(proj.ID)
		if !ok {
			return ErrProjectNotFound
		}
		if !p.DeleteUser(user.ID) {
			return ErrUserNotFound
		}
	}
	// Delete from users.
	if !c.users.Delete(user.ID) {
		return ErrUserNotFound
	}
	return nil
}
