// Copyright 2015 Eleme Inc. All rights reserved.

package persist

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// AddProject adds a proj to db.
func (p *Persist) AddProject(proj *models.Project) error {
	if err := p.db.Create(proj).Error; err != nil {
		switch err {
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

// UpdateProject updates a project with another.
func (p *Persist) UpdateProject(proj *models.Project) error {
	if err := p.db.Model(proj).Update(proj).Error; err != nil {
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

// DeleteProject deletes a project.
func (p *Persist) DeleteProject(proj *models.Project) error {
	if err := p.db.Delete(proj).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

// Projects returns all projects.
func (p *Persist) Projects(projs *[]*models.Project) error {
	var res []models.Project
	if err := p.db.Find(&res).Error; err != nil {
		return err
	}
	for _, proj := range res {
		*projs = append(*projs, &proj)
	}
	return nil
}

// RulesOfProject returns all rules for given project.
func (p *Persist) RulesOfProjects(proj *models.Project, rules *[]*models.Rule) error {
	var res []models.Rule
	if err := p.db.Model(proj).Related(res).Error; err != nil {
		return err
	}
	for _, rule := range res {
		*rules = append(*rules, &rule)
	}
	return nil
}

// UsersOfProjects returns all users for given project.
func (p *Project) UsersOfProjects(proj *models.Project, users *[]*models.User) error {
	var res []models.User
	if err := p.db.Model(proj).Related(res, "Users").Error; err != nil {
		return err
	}
	for _, user := range res {
		*users = append(*users, &user)
	}
	return nil
}
