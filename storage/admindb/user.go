// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// Users returns all the users.
func (db *DB) Users() (l []*models.User) {
	for _, user := range db.users.Items() {
		l = append(l, user.(*models.User))
	}
	return l
}

// GetUser returns user by id.
func (db *DB) GetUser(id int) (*models.User, error) {
	v, ok := db.users.Get(id)
	if !ok {
		return nil, ErrNotFound
	}
	return v.(*models.User), nil
}

// AddUser adds a user to db.
func (db *DB) AddUser(user *models.User) error {
	// Sql
	if err := db.db.Create(user).Error; err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		return err
	}
	// Cache
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
	u, err := db.GetUser(user.ID)
	if err != nil {
		return err
	}
	u.Update(user)
	// Update user in its projects.
	user = u.Clone()
	for _, p := range user.Projects {
		proj, err := db.GetProject(p.ID)
		if err != nil {
			return err
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
	user, err := db.GetUser(id)
	if err != nil {
		return err
	}
	// Clone
	user = user.Clone()
	// Delete user from its projects.
	for _, p := range user.Projects {
		proj, err := db.GetProject(p.ID)
		if err != nil {
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
