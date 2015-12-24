// Copyright 2015 Eleme Inc. All rights reserved.

package cache

import (
	// "github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage/admindb/persist"
	// "github.com/eleme/banshee/util/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// Open persist
	fileName := "db-testing"
	p, _ := persist.Open(fileName)
	// File should exist.
	defer p.Close()
	defer os.RemoveAll(fileName)
	// Init
	// TODO
}
