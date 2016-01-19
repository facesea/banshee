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
	rule1 := &models.Rule{Pattern: "a.b.*"}
	rule2 := &models.Rule{Pattern: "a.b.*.c"}
	rule3 := &models.Rule{Pattern: "a.*.c.d"}
	// Add to db.
	db.DB().Create(rule1)
	db.DB().Create(rule2)
	db.DB().Create(rule3)
	// Clear cache.
	db.RulesCache.rules.Clear()
	// Reload
	assert.Ok(t, nil == db.RulesCache.Init(db.DB()))
	// Get rule
	r1, ok := db.RulesCache.Get(rule1.ID)
	assert.Ok(t, ok)
	assert.Ok(t, r1.Pattern == rule1.Pattern)
	r2, ok := db.RulesCache.Get(rule2.ID)
	assert.Ok(t, ok)
	assert.Ok(t, r2.Pattern == rule2.Pattern)
	r3, ok := db.RulesCache.Get(rule3.ID)
	assert.Ok(t, ok)
	assert.Ok(t, r3.Pattern == rule3.Pattern)
}
