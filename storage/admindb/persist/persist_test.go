// Copyright 2015 Eleme Inc. All rights reserved.

package persist

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/assert"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	fileName := "db-testing"
	p, err := Open(fileName)
	// File should exist.
	assert.Ok(t, err == nil)
	assert.Ok(t, util.IsFileExist(fileName))
	defer p.Close()
	defer os.RemoveAll(fileName)
	// Tables should exist.
	assert.Ok(t, p.db.HasTable(&models.Project{}))
	assert.Ok(t, p.db.HasTable(&models.User{}))
	assert.Ok(t, p.db.HasTable(&models.Rule{}))
}
