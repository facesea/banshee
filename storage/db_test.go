// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/assert"
	"os"
	"path"
	"testing"
)

func TestOpen(t *testing.T) {
	// Open db.
	fileName := "storage_test"
	options := &Options{NumGrid: 288, GridLen: 300}
	db, err := Open(fileName, options)
	assert.Ok(t, err == nil)
	assert.Ok(t, db != nil)
	// Defer close and remove files.
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Check if child db file exist
	assert.Ok(t, util.IsFileExist(path.Join(fileName, admindbFileName)))
	assert.Ok(t, util.IsFileExist(path.Join(fileName, metricdbFileName)))
	sFileName := fmt.Sprintf("%s-%dx%d", statedbFileName, options.NumGrid, options.GridLen)
	assert.Ok(t, util.IsFileExist(path.Join(fileName, sFileName)))
}
