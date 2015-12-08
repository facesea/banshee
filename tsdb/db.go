// Package tsdb implements an embededable time series databas storage based on
// leveldb.

package tsdb

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// DB is a tsdb database.
type DB struct {
	db        *leveldb.DB
	tsLenLock *sync.Mutex
	hLenLock  *sync.Mutex
}

// Open a DB for given path, directory will be created if the path dose
// not exist. example:
//	db := OpenFile("mydb", nil)
//	defer db.Close()
func OpenFile(fileName string, options *opt.Options) (*DB, error) {
	ldb, err := leveldb.OpenFile(fileName, options)
	if err != nil {
		return nil, NewErrCorrupted(err)
	}
	db := new(DB)
	db.db = ldb
	db.tsLenLock = &sync.Mutex{}
	db.hLenLock = &sync.Mutex{}
	return db, nil
}

// Close a DB.
func (db *DB) Close() error {
	return db.db.Close()
}
