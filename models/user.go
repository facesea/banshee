// Copyright 2015 Eleme Inc. All rights reserved.

package models

// User is the alerter message receiver.
type User struct {
	// User may be cached.
	cache `sql:"-"`
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

// CopyIfShared returns a copy if the user is shared.
func (user *User) CopyIfShared() *User {
	if user.IsShared() {
		return user.Copy()
	}
	return user
}

// Copy the user.
func (user *User) Copy() *User {
	dst := &User{}
	user.CopyTo(dst)
	return dst
}

// CopyTo copy the user to another.
func (user *User) CopyTo(u *User) {
	user.RLockIfShared()
	defer user.RUnlockIfShared()
	u.LockIfShared()
	defer u.UnlockIfShared()
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
	u.EnableEmail = user.EnableEmail
	u.Phone = user.Phone
	u.EnablePhone = user.EnablePhone
	u.Projects = make([]*Project, len(user.Projects))
	copy(u.Projects, user.Projects)
}

// Equal tests the equality.
func (user *User) Equal(u *User) bool {
	user.RLockIfShared()
	defer user.RUnlockIfShared()
	u.LockIfShared()
	defer u.UnlockIfShared()
	if u.ID != user.ID {
		return false
	}
	if u.Name != user.Name {
		return false
	}
	if u.Email != user.Email {
		return false
	}
	if u.EnableEmail != user.EnableEmail {
		return false
	}
	if u.Phone != user.Phone {
		return false
	}
	if u.EnablePhone != user.EnablePhone {
		return false
	}
	// Projects
	if len(u.Projects) != len(user.Projects) {
		return false
	}
	for i := 0; i < len(u.Projects); i++ {
		if u.Projects[i].ID != user.Projects[i].ID {
			return false
		}
	}
	return true
}

// AddProject adds a project to the user.
func (user *User) AddProject(proj *Project) {
	proj = proj.CopyIfShared()
	user.LockIfShared()
	defer user.UnlockIfShared()
	user.Projects = append(user.Projects, proj)
}

// DeleteProject deletes a project from user.
func (user *User) DeleteProject(id int) bool {
	user.LockIfShared()
	defer user.UnlockIfShared()
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
	proj = proj.CopyIfShared()
	user.LockIfShared()
	defer user.UnlockIfShared()
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
	user.LockIfShared()
	defer user.UnlockIfShared()
	user.Name = u.Name
	user.Email = u.Email
	user.EnableEmail = u.EnableEmail
	user.Phone = u.Phone
	user.EnablePhone = u.EnablePhone
}

// GetProjects returns the projects of the user.
func (user *User) GetProjects() []*Project {
	user.RLockIfShared()
	defer user.RUnlockIfShared()
	if user.IsShared() {
		l := make([]*Project, len(user.Projects))
		copy(l, user.Projects)
		return l
	}
	return user.Projects
}

// GetProject returns the project of the user.
func (user *User) GetProject(id int) (*Project, bool) {
	user.RLockIfShared()
	defer user.RUnlockIfShared()
	for _, proj := range user.Projects {
		if proj.ID == id {
			return proj, true
		}
	}
	return nil, false
}
