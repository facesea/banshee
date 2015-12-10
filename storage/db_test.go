// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/eleme/banshee/config"
	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	cfg := config.NewConfigWithDefaults()
	cfg.Storage.Path = "storage_test/"
	db, err := Open(cfg)
	assert.Nil(t, err)
	defer db.Close()
	_, err = os.Stat(path.Join(cfg.Storage.Path, "rules"))
	assert.Nil(t, err)
	s := fmt.Sprintf("%dx%d", cfg.Periodicity[0], cfg.Periodicity[1])
	_, err = os.Stat(path.Join(cfg.Storage.Path, s))
	assert.Nil(t, err)
	assert.Nil(t, os.RemoveAll(cfg.Storage.Path))
}
