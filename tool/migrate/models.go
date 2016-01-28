// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"time"
)

// bell models (schema)
// https://github.com/eleme/bell.js/blob/master/lib/models.js

// Project => models.Project.
type Project struct {
	ID       int       `gorm:"primary_key;column:id"`
	Name     string    `gorm:"column:name" sql:"not null;unique"`
	Rules    []*Rule   `gorm:"foreignkey:ProjectId"`
	CreateAt time.Time `gorm:"column:createAt"`
	UpdateAt time.Time `gorm:"column:updateAt"`
}

// Rule => models.rule.
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

// Receiver => models.User.
type Receiver struct {
	ID          int       `gorm:"primary_key;column:id"`
	Name        string    `gorm:"column:name" sql:"index;not null;unique"`
	Email       string    `gorm:"column:email"`
	EnableEmail bool      `gorm:"column:enableEmail"`
	Phone       string    `gorm:"column:phone"`
	EnablePhone bool      `gorm:"column:enablePhone"`
	Universal   bool      `gorm:"column:universal`
	CreateAt    time.Time `gorm:"column:createAt"`
	UpdateAt    time.Time `gorm:"column:updateAt"`
}

// ReceiverProject is the join table.
//
// I have tried the struct tags, and patch myself, but all just not work. T_T
// Related issue: https://github.com/jinzhu/gorm/issues/707
// Related stackoverflow:
// http://stackoverflow.com/questions/33437453/using-different-join-table-columns-than-defaults-with-gorm-many2many-join
//
type ReceiverProject struct {
	ReceiverID int       `gorm:"primary_key;column:ReceiverId"`
	ProjectID  int       `gorm:"primary_key;column:ProjectId"`
	CreateAt   time.Time `gorm:"column:createAt"`
	UpdateAt   time.Time `gorm:"column:updateAt"`
}

// TableName for ReceiverProject.
func (rp *ReceiverProject) TableName() string {
	return "ReceiverProjects"
}
