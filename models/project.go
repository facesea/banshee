// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Project is a rules group.
type Project struct {
	// Project may be cached.
	cache `sql:"-" json:"-"`
	// ID in db.
	ID int `json:"-"`
	// Name
	Name string `sql:"not null";unique`
	// Project may have many rules, they shouldn't be shared.
	Rules []*Rule `json:"-"`
	// Project may have many users, they shouldn't be shared.
	Users []*User `gorm:"many2many:project_users" json:"-"`
}

// Copy the project.
func (proj *Project) Copy() *Project {
	if proj.IsShared() {
		// Lock if shared.
		proj.RLock()
		defer proj.RUnlock()
	}
	dst := &Project{ID: proj.ID, Name: proj.Name}
	copy(dst.Rules, proj.Rules)
	copy(dst.Users, proj.Users)
	return dst
}

// AddRule adds a rule to the project.
func (proj *Project) AddRule(rule *Rule) {
	if rule.IsShared() {
		// Copy if shared.
		rule = rule.Copy()
	}
	if proj.IsShared() {
		// Lock if shared.
		proj.Lock()
		defer proj.Unlock()
	}
	proj.Rules = append(proj.Rules, rule)
}

// AddUser adds a user to the project.
func (proj *Project) AddUser(user *User) {
	if user.IsShared() {
		// Copy if shared.
		user = user.Copy()
	}
	if proj.IsShared() {
		// Lock if shared.
		proj.Lock()
		defer proj.Unlock()
	}
	proj.Users = append(proj.Users, user)
}

// DeleteRule deletes a rule from the project.
func (proj *Project) DeleteRule(id int) bool {
	if proj.IsShared() {
		// Lock if shared.
		proj.Lock()
		defer proj.Unlock()
	}
	for i, rule := range proj.Rules {
		if rule.ID == id {
			proj.Rules = append(proj.Rules[:i], proj.Rules[i+1:]...)
			return true
		}
	}
	return false
}

// DeleteUser deletes a user from the project.
func (proj *Project) DeleteUser(id int) bool {
	if proj.IsShared() {
		// Lock if shared.
		proj.Lock()
		defer proj.Unlock()
	}
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
	if user.IsShared() {
		// Copy if shared.
		user = user.Copy()
	}
	if proj.IsShared() {
		// Lock if shared
		proj.Lock()
		defer proj.Unlock()
	}
	for i, u := range proj.Users {
		if u.ID == user.ID {
			tmp := u.Copy()
			tmp.Update(user)
			proj.Users[i] = tmp
			return true
		}
	}
	return false
}

// Update the project.
func (proj *Project) Update(project *Project) {
	if project.IsShared() {
		// Copy if shared.
		project = project.Copy()
	}
	if proj.IsShared() {
		// Lock if shared
		proj.Lock()
		defer proj.Unlock()
	}
	proj.Name = project.Name
}

// GetRules returns the rules of the project.
func (proj *Project) GetRules() []*Rule {
	if proj.IsShared() {
		// RLock if shared.
		proj.RLock()
		defer proj.RUnlock()
		var l []*Rule
		copy(l, proj.Rules)
		return l
	}
	// Return itself.
	return proj.Rules
}

// GetUsers returns the users of the project.
func (proj *Project) GetUsers() []*User {
	if proj.IsShared() {
		// RLock if shared.
		proj.RLock()
		defer proj.RUnlock()
		var l []*User
		copy(l, proj.Users)
		return l
	}
	// Return itself.
	return proj.Users
}
