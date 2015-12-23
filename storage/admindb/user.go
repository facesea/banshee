// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// getUser returns the user by id.
func (db *DB) getUser(id int) (*models.User, bool) {
	v, ok := db.users.Get(id)
	if !ok {
		// Not found.
		return nil, false
	}
	user := v.(*models.User)
	return user, true
}

// Users returns all the users.
func (db *DB) Users() (l []*models.User) {
	for _, v := range db.users.Items() {
		user := v.(*models.User)
		l = append(l, user.Copy())
	}
	return l
}

// GetUser returns user by into a local value.
func (db *DB) GetUser(u *models.User) error {
	user, ok := db.getUser(u.ID)
	if !ok {
		return ErrNotFound
	}
	user.CopyTo(u)
	return nil
}

// AddUser adds a user to db.
func (db *DB) AddUser(user *models.User) error {
	// Sql: user.ID will be created.
	if err := db.db.Create(user).Error; err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		if err == sqlite3.ErrConstraintNotNull {
			return ErrConstraintNotNull
		}
		return err
	}
	// Cache a copy.
	user = user.Copy()
	// Mark as shared.
	user.MakeShared()
	// Add to users
	db.users.Set(user.ID, user)
	return nil
}

// UpdateUser updates a user with another.
func (db *DB) UpdateUser(user *models.User) error {
	// Sql
	if err := db.db.Model(user).Update(user).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		return err
	}
	// Cache
	// Update user in users.
	u, ok := db.getUser(user.ID)
	if !ok {
		return ErrNotFound
	}
	u.Update(user)
	// Update user in its projects.
	projects := u.GetProjects()
	for _, p := range projects {
		proj, ok := db.getProject(p.ID)
		if !ok {
			return ErrNotFound
		}
		if !proj.UpdateUser(user) {
			return ErrNotFound
		}
	}
	return nil
}

// DeleteUser deletes a user from db.
func (db *DB) DeleteUser(id int) error {
	// Sql
	if err := db.db.Delete(&models.User{ID: id}).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		return err
	}
	// Cache
	// Get this user.
	user, ok := db.getUser(id)
	if !ok {
		return ErrNotFound
	}
	// Delete user from its projects.
	projects := user.GetProjects()
	for _, p := range projects {
		proj, ok := db.getProject(p.ID)
		if !ok {
			return ErrNotFound
		}
		if !proj.DeleteUser(id) {
			return ErrNotFound
		}
	}
	// Delete user from users.
	if !db.users.Delete(id) {
		return ErrNotFound
	}
	return nil
}
