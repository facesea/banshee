// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// Projects returns all the projects.
func (db *DB) Projects() (l []*models.Project) {
	for _, proj := range db.projects.Items() {
		l = append(l, proj.(*models.Project))
	}
	return l
}

// GetProject returns project by id.
func (db *DB) GetProject(id int) (*models.Project, error) {
	v, ok := db.projects.Get(id)
	if !ok {
		return nil, ErrNotFound
	}
	return v.(*models.Project), nil
}

// AddProject adds a project to db.
func (db *DB) AddProject(proj *models.Project) error {
	// Sql
	if err := db.db.Create(proj).Error; err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		return err
	}
	// Cache
	// Add to projects.
	db.projects.Set(proj.ID, proj)
	return nil
}

// UpdateProject updates a project with another.
func (db *DB) UpdateProject(proj *models.Project) error {
	// Sql
	if err := db.db.Model(proj).Update(proj).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		if err == sqlite3.ErrConstraintUnique {
			return ErrConstraintUnique
		}
		return err
	}
	// Cache
	// Update project in projects.
	project, err := db.GetProject(proj.ID)
	if err != nil {
		return err
	}
	project.Update(proj)
	// Update project in its users.
	proj = project.Clone()
	for _, u := range project.Users {
		user, err := db.GetUser(u.ID)
		if err != nil {
			return err
		}
		if !user.UpdateProject(proj) {
			return ErrNotFound
		}
	}
	return nil
}

// DeleteProject deletes a project from db.
func (db *DB) DeleteProject(id int) error {
	// Sql
	if err := db.db.Delete(&models.Project{ID: id}).Error; err != nil {
		if err == gorm.RecordNotFound {
			return ErrNotFound
		}
		return err
	}
	// Cache
	// Get this project.
	proj, err := db.GetProject(id)
	if err != nil {
		return err
	}
	// Clone
	proj = proj.Clone()
	// Delete its rules.
	for _, rule := range proj.Rules {
		if !db.rules.Delete(rule.ID) {
			return ErrNotFound
		}
	}
	// Delete project from its users.
	for _, u := range proj.Users {
		user, err := db.GetUser(u.ID)
		if err != nil {
			return ErrNotFound
		}
		if !user.DeleteProject(id) {
			return ErrNotFound
		}
	}
	// Delete project from projects.
	if !db.projects.Delete(id) {
		return ErrNotFound
	}
	return nil
}
