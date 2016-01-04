// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	rule := &models.Rule{Pattern: "a.b.*"}
	// Add one to db.
	db.DB().Create(rule)
	// Clear cache.
	db.RulesCache.rules.Clear()
	// Reload
	assert.Ok(t, nil == db.RulesCache.Init(db.DB()))
	// Get rule
	r, ok := db.RulesCache.Get(rule.ID)
	assert.Ok(t, ok)
	assert.Ok(t, r.Pattern == rule.Pattern)
}
