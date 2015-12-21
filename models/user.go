// Copyright 2015 Eleme Inc. All rights reserved.

package models

import "sync"

// User is the alerter message receiver.
// One user can receive many projects.
type User struct {
	// A user instance will be shared between goroutines as cache, this RWMetux
	// is to guarantee its safety.
	sync.RWMutex
	// ID in db.
	ID int `gorm:"primary_key"`
	// Name
	Name string `sql:"not null;unique"`
	// Email
	Email       string
	EnableEmail bool
	// Phone
	Phone       string
	EnablePhone bool
	// Users can subscribe many projects.
	Projects []*Project `gorm:"many2many:project_users"`
}

// AddProject adds a project to the user.
func (user *User) AddProject(proj *Project) {
	user.Lock()
	defer user.Unlock()
	user.Projects = append(user.Projects, proj)
}

// DeleteProject deletes a project from user.
func (user *User) DeleteProject(id int) bool {
	user.Lock()
	defer user.Unlock()
	for i, proj := range user.Projects {
		if proj.ID == id {
			user.Projects = append(user.Projects[:i], user.Projects[i+1:]...)
			return true
		}
	}
	return false
}

// UpdateProject updates a project.
func (user *User) UpdateProject(proj *Project) bool {
	user.Lock()
	defer user.Unlock()
	for i, p := range user.Projects {
		if p.ID == proj.ID {
			user.Projects[i] = proj
			return true
		}
	}
	return false
}

// Update the user.
func (user *User) Update(u *User) {
	user.Lock()
	defer user.Unlock()
	user.Name = u.Name
	user.Email = u.Email
	user.EnableEmail = u.EnableEmail
	user.Phone = u.Phone
	user.EnablePhone = u.EnablePhone
}

// Clone the user.
func (user *User) Clone() *User {
	user.RLock()
	defer user.RUnlock()
	dst := &User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		EnableEmail: user.EnableEmail,
		Phone:       user.Phone,
		EnablePhone: user.EnablePhone,
	}
	copy(dst.Projects, user.Projects)
	return dst
}
