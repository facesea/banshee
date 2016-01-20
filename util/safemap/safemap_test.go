// Copyright 2015 Eleme Inc. All rights reserved.

package safemap

import (
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestBasic(t *testing.T) {
	m := New()
	// Set
	m.Set("key1", "val1")
	m.Set("key2", "val2")
	m.Set("key3", "val3")
	assert.Ok(t, m.Len() == 3)
	// Get
	val1, ok := m.Get("key1")
	assert.Ok(t, ok)
	assert.Ok(t, val1 == "val1")
	// Items
	assert.Ok(t, m.Items()["key1"] == "val1")
	assert.Ok(t, m.Items()["key2"] == "val2")
	assert.Ok(t, m.Items()["key3"] == "val3")
	// Len
	assert.Ok(t, m.Len() == 3)
	// Delete
	assert.Ok(t, m.Delete("key1"))
	assert.Ok(t, !m.Delete("key-not-exist"))
	assert.Ok(t, m.Len() == 2)
	// Clear
	m.Clear()
	assert.Ok(t, m.Len() == 0)
}
