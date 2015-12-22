// Copyright 2015 Eleme Inc. All rights reserved.

package models

// User is the alerter message receiver.
type User struct {
	// User may be cached.
	cache `sql:"-" json:"-"`
	// ID in db.
	ID int `gorm:"primary_key" json:"-"`
	// Name
	Name string `sql:"not null;unique"`
	// Email
	Email       string
	EnableEmail bool
	// Phone
	Phone       string
	EnablePhone bool
	// Users can subscribe many projects.
	Projects []*Project `gorm:"many2many:project_users" json:"-"`
}

// Copy the user.
func (user *User) Copy() *User {
	if user.IsShared() {
		user.RLock()
		defer user.RUnlock()
	}
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

// AddProject adds a project to the user.
func (user *User) AddProject(proj *Project) {
	if proj.IsShared() {
		// Copy if shared.
		proj = proj.Copy()
	}
	if user.IsShared() {
		// Lock if shared.
		user.Lock()
		defer user.Unlock()
	}
	user.Projects = append(user.Projects, proj)
}

// DeleteProject deletes a project from user.
func (user *User) DeleteProject(id int) bool {
	if user.IsShared() {
		// Lock if shared.
		user.Lock()
		defer user.Unlock()
	}
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
	if proj.IsShared() {
		// Copy if shared.
		proj = proj.Copy()
	}
	if user.IsShared() {
		// Lock if shared.
		user.Lock()
		defer user.Unlock()
	}
	for i, p := range user.Projects {
		if p.ID == proj.ID {
			tmp := p.Copy()
			tmp.Update(proj)
			user.Projects[i] = tmp
			return true
		}
	}
	return false
}

// Update the user.
func (user *User) Update(u *User) {
	if user.IsShared() {
		// Lock if shared.
		user.Lock()
		defer user.Unlock()
	}
	user.Name = u.Name
	user.Email = u.Email
	user.EnableEmail = u.EnableEmail
	user.Phone = u.Phone
	user.EnablePhone = u.EnablePhone
}

// GetProjects returns the projects of the user.
func (user *User) GetProjects() []*Project {
	if user.IsShared() {
		// RLock if shared.
		user.RLock()
		defer user.RUnlock()
		// Return a copy.
		var l []*Project
		copy(l, user.Projects)
		return l
	}
	// Return itself.
	return user.Projects
}
