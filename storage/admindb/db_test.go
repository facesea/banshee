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
	assert.Ok(t, nil == err)
	assert.Ok(t, db != nil)
	assert.Ok(t, util.IsFileExist(fileName))
	defer os.RemoveAll(fileName)
	defer db.Close()
	assert.Ok(t, db.DB().HasTable(&models.User{}))
	assert.Ok(t, db.DB().HasTable(&models.Rule{}))
	assert.Ok(t, db.DB().HasTable(&models.Project{}))
}
