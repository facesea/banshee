// Copyright 2015 Eleme Inc. All rights reserved.

package persist

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

// AddRule adds a rule to db.
func (p *Persist) AddRule(rule *models.Rule) error {
	if err := p.db.Create(rule).Error; err != nil {
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

// DeleteRule deletes a rule from db.
func (p *Persist) DeleteRule(rule *models.Rule) error {
	if err := p.db.Delete(rule).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

// GetRules returns all rules.
func (p *Persist) GetRules(rules *[]*models.Rule) error {
	var res []models.Rule
	if err := p.db.Find(&res).Error; err != nil {
		return err
	}
	for _, rule := range res {
		*rules = append(*rules, &rule)
	}
	return nil
}
