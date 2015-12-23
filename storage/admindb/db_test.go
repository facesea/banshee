// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/assert"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	fileName := "db-testing"
	db, err := Open(fileName)
	// File should exist.
	assert.Ok(t, err == nil)
	assert.Ok(t, util.IsFileExist(fileName))
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Tables should exist.
	assert.Ok(t, db.db.HasTable(&models.Project{}))
	assert.Ok(t, db.db.HasTable(&models.User{}))
	assert.Ok(t, db.db.HasTable(&models.Rule{}))
}

func TestLoad(t *testing.T) {
	fileName := "db-testing"
	db, err := Open(fileName)
	// File should exist.
	assert.Ok(t, err == nil)
	assert.Ok(t, util.IsFileExist(fileName))
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add user.
	user := &models.User{Name: "hit9"}
	assert.Ok(t, nil == db.AddUser(user))
	// Add proj.
	proj := &models.Project{Name: "banshee"}
	assert.Ok(t, nil == db.AddProject(proj))
	// Add rule.
	rule := &models.Rule{ProjectID: proj.ID, Pattern: "counter.*"}
	assert.Ok(t, nil == db.AddRuleToProject(proj, rule))
	// Add user to proj.
	assert.Ok(t, nil == db.AddUserToProject(proj, user))
	// Clear the cache.
	db.users.Clear()
	db.projects.Clear()
	db.rules.Clear()
	// Reload the cache.
	assert.Ok(t, nil == db.load())
	// Check cache.
	// Must be not empty.
	assert.Ok(t, db.users.Len() == 1)
	assert.Ok(t, db.rules.Len() == 1)
	assert.Ok(t, db.projects.Len() == 1)
	// Reloaded rule should be equal with old rule.
	r := &models.Rule{ID: rule.ID}
	assert.Ok(t, nil == db.GetRule(r))
	assert.Ok(t, r.Equal(rule))
	// Reloaded proj should be equal with old proj.
	p := &models.Project{ID: proj.ID}
	assert.Ok(t, nil == db.GetProject(p))
	assert.Ok(t, p.Equal(proj))
	// Reloaded user should be equal with old user.
	u := &models.User{ID: user.ID}
	assert.Ok(t, nil == db.GetUser(u))
	assert.Ok(t, u.Equal(user))
}
