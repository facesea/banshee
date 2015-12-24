// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import "github.com/eleme/banshee/models"

// NumUsers returns the number of users.
func (db *DB) NumUsers() int {
	return db.cache.NumUsers()
}

// GetUser returns user.
func (db *DB) GetUser(user *models.User) error {
	return db.cache.GetUser(user)
}

// Users returns all users.
func (db *DB) Users(users *[]*models.User) {
	db.cache.Users(users)
}

// UsersN returns users for given range.
func (db *DB) UsersN(users *[]*models.User, offset int, limit int) {
	db.cache.UsersN(users, offset, limit)
}

// AddUser adds a user to db.
func (db *DB) AddUser(user *models.User) error {
	if err := db.persist.AddUser(user); err != nil {
		return err
	}
	db.cache.AddUser(user)
	return nil
}

// UpdateUser updates a user with another.
func (db *DB) UpdateUser(user *models.User) error {
	if err := db.persist.UpdateUser(user); err != nil {
		return err
	}
	return db.cache.UpdateUser(user)
}

// DeleteUser deletes a user from db.
func (db *DB) DeleteUser(user *models.User) error {
	if err := db.persist.DeleteUser(user); err != nil {
		return err
	}
	return db.cache.DeleteUser(user)
}
