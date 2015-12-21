// Copyright 2015 Eleme Inc. All rights reserved.

package models

import "sync"

// Project is rules group.
// One project can have many rules and many recievers.
type Project struct {
	// A project instance will be shared between goroutines as cache, this
	// RWMutex is to guarantee its safety.
	sync.RWMutex
	// ID in db.
	ID int `gorm:"primary_key"`
	// Name
	Name string `sql:"not null;unique"`
	// Project has many rules.
	Rules []*Rule
	// Project has many receivers.
	Users []*User `gorm:"many2many:project_user"`
}

// AddRule adds a rule to the project.
func (proj *Project) AddRule(rule *Rule) {
	proj.Lock()
	defer proj.Unlock()
	proj.Rules = append(proj.Rules, rule)
}

// AddUser adds a user to the project.
func (proj *Project) AddUser(user *User) {
	proj.Lock()
	defer proj.Unlock()
	proj.Users = append(proj.Users, user)
}

// DeleteRule deletes a rule from project.
func (proj *Project) DeleteRule(id int) bool {
	proj.Lock()
	defer proj.Unlock()
	for i, rule := range proj.Rules {
		if rule.ID == id {
			proj.Rules = append(proj.Rules[:i], proj.Rules[i+1:]...)
			return true
		}
	}
	return false
}

// DeleteUser deletes a user from project.
func (proj *Project) DeleteUser(id int) bool {
	proj.Lock()
	defer proj.Unlock()
	for i, user := range proj.Users {
		if user.ID == id {
			proj.Users = append(proj.Users[:i], proj.Users[i+1:]...)
			return true
		}
	}
	return false
}

// UpdateUser updates a user.
func (proj *Project) UpdateUser(user *User) bool {
	proj.Lock()
	defer proj.Unlock()
	for i, u := range proj.Users {
		if u.ID == user.ID {
			proj.Users[i] = user
			return true
		}
	}
	return false
}

// Update updates the project.
func (proj *Project) Update(project *Project) {
	proj.Lock()
	defer proj.Unlock()
	proj.Name = project.Name
}

// Clone the project.
func (proj *Project) Clone() *Project {
	proj.RLock()
	defer proj.RUnlock()
	dst := &Project{
		ID:   proj.ID,
		Name: proj.Name,
	}
	copy(dst.Rules, proj.Rules)
	copy(dst.Users, proj.Users)
	return dst
}
