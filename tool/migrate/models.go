// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

// bell models (schema)
// https://github.com/eleme/bell.js/blob/master/lib/models.js

type Project struct {
	ID        int         `gorm:"primary_key;column:id"`
	Name      string      `gorm:"column:name" sql:"not null;unique"`
	Rules     []*Rule     `gorm:"foreignkey:ProjectId"`
	Receivers []*Receiver `gorm:"many2many:ReceiverProjects"`
	CreateAt  time.Time   `gorm:"column:createAt"`
	UpdateAt  time.Time   `gorm:"column:updateAt"`
}

type Rule struct {
	ID        int       `gorm:"primary_key;column:id"`
	Pattern   string    `gorm:"column:pattern" sql:"unique"`
	Up        bool      `gorm:"column:up"`
	Down      bool      `gorm:"column:down"`
	Min       float64   `gorm:"column:min"`
	Max       float64   `gorm:"column:max"`
	ProjectID int       `gorm:"column:ProjectId;associationforeignkey:ProjectId"`
	CreateAt  time.Time `gorm:"column:createAt"`
	UpdateAt  time.Time `gorm:"column:updateAt"`
}

type Receiver struct {
	ID          int        `gorm:"primary_key;column:id"`
	Name        string     `gorm:"column:name" sql:"index;not null;unique"`
	Email       string     `gorm:"column:email"`
	EnableEmail bool       `gorm:"column:enableEmail"`
	Phone       string     `gorm:"column:phone"`
	EnablePhone bool       `gorm:"column:enablePhone"`
	Universal   bool       `gorm:"column:universal`
	Projects    []*Project `gorm:"many2many:ReceiverProjects"`
	CreateAt    time.Time  `gorm:"column:createAt"`
	UpdateAt    time.Time  `gorm:"column:updateAt"`
}

// Patch gorm to change the join table column names for `ReceiverProjects`.
// Related issue: https://github.com/jinzhu/gorm/issues/707
func patchReceiverProjectsFieldNames(db *gorm.DB) {
	for _, field := range db.NewScope(&Receiver{}).GetStructFields() {
		if field.Name == "Projects" {
			field.Relationship.ForeignDBNames = []string{"ReceiverId"}
			field.Relationship.AssociationForeignDBNames = []string{"ProjectId"}
		}
	}
}
