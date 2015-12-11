// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/util"
	"os"
	"path"
	"testing"
)

func TestOpen(t *testing.T) {
	cfg := config.NewConfigWithDefaults()
	cfg.Storage.Path = "storage_test/"
	db, err := Open(cfg)
	util.Assert(t, err == nil)
	defer db.Close()
	_, err = os.Stat(path.Join(cfg.Storage.Path, "rules"))
	util.Assert(t, err == nil)
	s := fmt.Sprintf("%dx%d", cfg.Periodicity[0], cfg.Periodicity[1])
	_, err = os.Stat(path.Join(cfg.Storage.Path, s))
	util.Assert(t, err == nil)
	util.Assert(t, os.RemoveAll(cfg.Storage.Path) == nil)
}
