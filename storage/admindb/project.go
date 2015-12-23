// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// getProject returns the project by id.
func (db *DB) getProject(id int) (*models.Project, bool) {
	v, ok := db.projects.Get(id)
	if !ok {
		// Not found.
		return nil, false
	}
	proj := v.(*models.Project)
	return proj, true
}

// Projects returns all the projects.
func (db *DB) Projects() (l []*models.Project) {
	for _, v := range db.projects.Items() {
		proj := v.(*models.Project)
		l = append(l, proj.Copy())
	}
	return l
}

// GetProject returns project into a local value.
func (db *DB) GetProject(p *models.Project) error {
	proj, ok := db.getProject(p.ID)
	if !ok {
		return ErrProjectNotFound
	}
	proj.CopyTo(p)
	return nil
}

// AddProject adds a project to db.
func (db *DB) AddProject(proj *models.Project) error {
	// Sql: proj.ID will be created.
	if err := db.db.Create(proj).Error; err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		if err == sqlite3.ErrConstraintNotNull {
			return ErrConstraintNotNull
		}
		return err
	}
	// Cache a copy.
	proj = proj.Copy()
	// Mark as shared.
	proj.MakeShared()
	// Add to projects.
	db.projects.Set(proj.ID, proj)
	return nil
}

// UpdateProject updates a project with another.
func (db *DB) UpdateProject(proj *models.Project) error {
	// Sql
	if err := db.db.Model(proj).Update(proj).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrProjectNotFound
		}
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		return err
	}
	// Cache
	// Update project in projects.
	project, ok := db.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	project.Update(proj)
	// Update project in its users.
	users := project.GetUsers()
	for _, u := range users {
		user, ok := db.getUser(u.ID)
		if !ok {
			return ErrUserNotFound
		}
		if !user.UpdateProject(proj) {
			return ErrProjectNotFound
		}
	}
	return nil
}

// DeleteProject deletes a project from db.
func (db *DB) DeleteProject(id int) error {
	// Sql
	if err := db.db.Delete(&models.Project{ID: id}).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrProjectNotFound
		}
		return err
	}
	// Cache
	// Get this project.
	proj, ok := db.getProject(id)
	if !ok {
		return ErrProjectNotFound
	}
	// Delete its rules.
	rules := proj.GetRules()
	for _, rule := range rules {
		if !db.rules.Delete(rule.ID) {
			return ErrRuleNotFound
		}
	}
	// Delete project from its users.
	users := proj.GetUsers()
	for _, u := range users {
		user, ok := db.getUser(u.ID)
		if !ok {
			return ErrUserNotFound
		}
		if !user.DeleteProject(id) {
			return ErrProjectNotFound
		}
	}
	// Delete project from projects.
	if !db.projects.Delete(id) {
		return ErrProjectNotFound
	}
	return nil
}

// AddUserToProject adds a user to a project.
func (db *DB) AddUserToProject(proj *models.Project, user *models.User) error {
	// If user exist
	u, ok := db.getUser(user.ID)
	if !ok {
		return ErrUserNotFound
	}
	// If proj exist
	p, ok := db.getProject(proj.ID)
	if !ok {
		return ErrProjectNotFound
	}
	// Sql: user will be appened to proj.Users.
	if err := db.db.Model(proj).Association("Users").Append(user).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		if err == sqlite3.ErrConstraintPrimaryKey {
			return ErrConstraintPrimaryKey
		}
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		return err
	}
	// Append proj to user.
	user.AddProject(proj)
	// Add user to project.
	p.AddUser(u)
	// Add project to user.
	u.AddProject(p)
	return nil
}
