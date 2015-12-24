// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import "github.com/eleme/banshee/models"

// NumProjects returns the number of projects.
func (db *DB) NumProjects() int {
	return db.cache.NumProjects()
}

// GetProjects returns all projects.
func (db *DB) GetProjects(projs *[]*models.Project) {
	db.cache.GetProjects(projs)
}

// GetProjectsN returns projects for given range.
func (db *DB) GetProjectsN(projs *[]*models.Project, offset int, limit int) {
	db.cache.GetProjectsN(projs, offset, limit)
}

// GetProject returns project.
func (db *DB) GetProject(proj *models.Project) error {
	return db.cache.GetProject(proj)
}

// AddProject adds a project to db.
func (db *DB) AddProject(proj *models.Project) error {
	if err := db.persist.AddProject(proj); err != nil {
		return err
	}
	db.cache.AddProject(proj)
	return nil
}

// UpdateProject updates a project with another.
func (db *DB) UpdateProject(proj *models.Project) error {
	if err := db.persist.UpdateProject(proj); err != nil {
		return err
	}
	return db.cache.UpdateProject(proj)
}

// DeleteProject deletes a project from db.
func (db *DB) DeleteProject(proj *models.Project) error {
	if err := db.persist.DeleteProject(proj); err != nil {
		return err
	}
	return db.cache.DeleteProject(proj)
}

// AddUserToProject adds a user to a project.
func (db *DB) AddUserToProject(proj *models.Project, user *models.User) error {
	if err := db.persist.AddUserToProject(proj, user); err != nil {
		return err
	}
	return db.cache.AddUserToProject(proj, user)
}
