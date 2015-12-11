// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"github.com/eleme/banshee/util"
	"os"
	"path"
	"testing"
)

func TestOpen(t *testing.T) {
	fileName := "./storage_test"
	numGrids, gridLen := 480, 180
	db, err := Open(fileName, numGrids, gridLen)
	util.Assert(t, err == nil)
	defer db.Close()
	_, err = os.Stat(path.Join(fileName, "rules"))
	util.Assert(t, err == nil)
	s := fmt.Sprintf("%dx%d", numGrids, gridLen)
	_, err = os.Stat(path.Join(fileName, s))
	util.Assert(t, err == nil)
	util.Assert(t, os.RemoveAll(fileName) == nil)
}
