// Copyright 2015 Eleme Inc. All rights reserved.

// Package storage implements a database for rules and detection states
// storage, based on leveldb.
//   db, err := Open(cfg)
//   if err != nil {
//       logger.Fatal("failed to open db: %v", err)
//   }
//   defer db.Close()
//
package storage

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/eleme/banshee/config"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

// DB file mode
const filemode = 0755

// A DB is a database.
type DB struct {
	// Detection
	d  *leveldb.DB
	dl *sync.Mutex
	// Rules
	r *leveldb.DB
	// Config
	cfg *config.Config
}

// Open a DB instance from config. If the storage.path dosen't exist, it will
// be created.
//
// Function Open will create 2 leveldb DB instance, one for rules storage and
// one for detection metrics storage.
//   storage/
//     |- rules
//     |- mxn
//
func Open(cfg *config.Config) (*DB, error) {
	p := cfg.Storage.Path
	// Try to make directory
	err := os.Mkdir(p, filemode)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	s := fmt.Sprintf("%dx%d", cfg.Periodicity[0], cfg.Periodicity[1])
	d, err := leveldb.OpenFile(path.Join(p, s), nil)
	if err != nil {
		return nil, err
	}
	r, err := leveldb.OpenFile(path.Join(p, "rules"), nil)
	if err != nil {
		return nil, err
	}
	db := new(DB)
	db.d = d
	db.r = r
	db.dl = &sync.Mutex{}
	db.cfg = cfg
	return db, nil
}

// Close a DB instance, this will close the rules db and metrics db.
func (db *DB) Close() {
	db.d.Close()
	db.r.Close()
}

// Help to test if the current error indicates the db corrupted.
func (db *DB) IsCorrupted(err error) bool {
	return leveldbErrors.IsCorrupted(err)
}
