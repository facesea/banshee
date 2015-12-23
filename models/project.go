// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Project is a rules group.
type Project struct {
	// Project may be cached.
	cache `sql:"-" json:"-"`
	// ID in db.
	ID int `json:"-"`
	// Name
	Name string `sql:"not null";unique json:"name"`
	// Project may have many rules, they shouldn't be shared.
	Rules []*Rule `json:"-"`
	// Project may have many users, they shouldn't be shared.
	Users []*User `gorm:"many2many:project_users" json:"-"`
}

// CopyIfShared returns a copy if the proj is shared.
func (proj *Project) CopyIfShared() *Project {
	if proj.IsShared() {
		return proj.Copy()
	}
	return proj
}

// Copy the project.
func (proj *Project) Copy() *Project {
	dst := &Project{}
	proj.CopyTo(dst)
	return dst
}

// CopyTo copy the project to another.
func (proj *Project) CopyTo(p *Project) {
	proj.RLockIfShared()
	defer proj.RUnlockIfShared()
	p.LockIfShared()
	defer p.UnlockIfShared()
	p.ID = proj.ID
	p.Name = proj.Name
	p.Rules = make([]*Rule, len(proj.Rules))
	p.Users = make([]*User, len(proj.Users))
	copy(p.Rules, proj.Rules)
	copy(p.Users, proj.Users)
}

// Equal tests the equality.
func (proj *Project) Equal(p *Project) bool {
	proj.RLockIfShared()
	defer proj.RLockIfShared()
	p.RLockIfShared()
	defer p.RLockIfShared()
	if p.ID != proj.ID {
		return false
	}
	if p.Name != proj.Name {
		return false
	}
	// Rules
	if len(p.Rules) != len(proj.Rules) {
		return false
	}
	for i := 0; i < len(p.Rules); i++ {
		found := false
		for j := 0; j < len(proj.Rules); j++ {
			if p.Rules[i].Equal(proj.Rules[j]) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	// Users
	if len(p.Users) != len(proj.Users) {
		return false
	}
	for i := 0; i < len(p.Users); i++ {
		found := false
		for j := 0; j < len(proj.Users); j++ {
			if p.Users[i].ID == proj.Users[j].ID {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// AddRule adds a rule to the project.
func (proj *Project) AddRule(rule *Rule) {
	rule = rule.CopyIfShared()
	proj.LockIfShared()
	defer proj.UnlockIfShared()
	proj.Rules = append(proj.Rules, rule)
}

// AddUser adds a user to the project.
func (proj *Project) AddUser(user *User) {
	user = user.CopyIfShared()
	proj.LockIfShared()
	defer proj.UnlockIfShared()
	proj.Users = append(proj.Users, user)
}

// DeleteRule deletes a rule from the project.
func (proj *Project) DeleteRule(id int) bool {
	proj.LockIfShared()
	defer proj.UnlockIfShared()
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
	proj.LockIfShared()
	defer proj.UnlockIfShared()
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
	user = user.CopyIfShared()
	proj.LockIfShared()
	defer proj.UnlockIfShared()
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
func (proj *Project) Update(p *Project) {
	p = p.CopyIfShared()
	proj.LockIfShared()
	defer proj.UnlockIfShared()
	proj.Name = p.Name
}

// GetRules returns the rules of the project.
func (proj *Project) GetRules() []*Rule {
	proj.RLockIfShared()
	defer proj.RUnlockIfShared()
	if proj.IsShared() {
		l := make([]*Rule, len(proj.Rules))
		copy(l, proj.Rules)
		return l
	}
	return proj.Rules
}

// GetUsers returns the users of the project.
func (proj *Project) GetUsers() []*User {
	proj.RLockIfShared()
	defer proj.RUnlockIfShared()
	if proj.IsShared() {
		l := make([]*User, len(proj.Users))
		copy(l, proj.Users)
		return l
	}
	return proj.Users
}

// GetUser returns a user of the project.
func (proj *Project) GetUser(id int) (*User, bool) {
	proj.RLockIfShared()
	defer proj.RUnlockIfShared()
	for _, user := range proj.Users {
		if user.ID == id {
			return user, true
		}
	}
	return nil, false
}
