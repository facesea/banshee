// Copyright 2015 Eleme Inc. All rights reserved.

package persist

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// AddUser adds a user to db.
func (p *Persist) AddUser(user *models.User) error {
	if err := p.db.Create(user).Error; err != nil {
		switch err {
		case sqlite3.ErrConstraintUnique:
			return ErrUnique
		case sqlite3.ErrConstraintNotNull:
			return ErrUnique
		case sqlite3.ErrConstraintPrimaryKey:
			return ErrPrimaryKey
		default:
			return err
		}
	}
	return nil
}

// UpdateUser updates a user with another.
func (p *Persist) UpdateUser(user *models.User) error {
	if err := p.db.Model(user).Update(user).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			return ErrNotFound
		case sqlite3.ErrConstraintNotNull:
			return ErrNotNull
		case sqlite3.ErrConstraintUnique:
			return ErrUnique
		case sqlite3.ErrConstraintPrimaryKey:
			return ErrPrimaryKey
		default:
			return err
		}
	}
	return nil
}

// DeleteUser deletes a user.
func (p *Persist) DeleteUser(user *models.User) error {
	if err := p.db.Delete(user).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

// Users return all users.
func (p *Persist) Users(users *[]*models.User) error {
	var res []models.User
	if err := p.db.Find(&res).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			return ErrNotFound
		default:
			return err
		}
	}
	for _, user := range res {
		*users = append(*users, &user)
	}
	return nil
}

// ProjectsOfUser return all projects for given user.
func (p *Persist) ProjectsOfUser(user *models.User, projs *[]*models.Project) error {
	var res []models.Project
	if err := p.db.Model(user).Related(&res, "Projects").Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			return ErrNotFound
		default:
			return err
		}
	}
	for _, proj := range res {
		*projs = append(*projs, &proj)
	}
	return nil
}
