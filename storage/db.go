// Copyright 2015 Eleme Inc. All rights reserved.

// Package storage implements a database for rules and detection states
// storage, based on leveldb.
//   db, err := Open("mydb", 480, 180)
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

	"github.com/syndtr/goleveldb/leveldb/errors"
)

// DB file mode
const filemode = 0755

// A DB is a database.
type DB struct {
	s *statesDB
	r *rulesDB
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
func Open(fileName string, numGrids, gridLen int) (*DB, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err = os.Mkdir(fileName, filemode)
		if err != nil {
			return nil, err
		}
	}
	name := path.Join(fileName, fmt.Sprintf("%dx%d", numGrids, gridLen))
	s, err := openStatesDB(name, numGrids, gridLen)
	if err != nil {
		return nil, err
	}
	name = path.Join(fileName, rulesFileName)
	r, err := openRulesDB(name)
	if err != nil {
		s.close()
		return nil, err
	}
	return &DB{s: s, r: r}, nil
}

// Close a DB instance, this will close the rules db and metrics db.
func (db *DB) Close() error {
	err := db.s.close()
	if err != nil {
		return err
	}
	err = db.r.close()
	if err != nil {
		return err
	}
	return nil
}

// Help to test if the current error indicates the db corrupted.
func isCorrupted(err error) bool {
	return errors.IsCorrupted(err)
}
