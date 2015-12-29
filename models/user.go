// Copyright 2015 Eleme Inc. All rights reserved.

package models

// User is the alerter message receiver.
type User struct {
	// ID in db.
	ID int `gorm:"primary_key" json:"id"`
	// Name
	Name string `sql:"not null;unique" json:"name"`
	// Email
	Email       string `json:"email"`
	EnableEmail bool   `json:"enableEmail"`
	// Phone
	Phone       string `json:"phone"`
	EnablePhone bool   `json:"enablePhone"`
	// Universal
	Universal bool `sql:"index" json:"universal"`
	// Users can subscribe many projects.
	Projects []*Project `gorm:"many2many:project_users" json:"-"`
}
