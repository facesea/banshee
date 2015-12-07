// Package tsdb implements an embededable time-series and hash-table
// storage based on leveldb.

package tsdb

import (
	"github.com/syndtr/goleveldb/leveldb"
	leveldbOpt "github.com/syndtr/goleveldb/leveldb/opt"
)

// A TimeStamp is auint64 number.
type TimeStamp uint64

// A Key is a string.
type Key string

// A Value is a float64 number.
type Value float64

// DB is a tsdb database.
type DB struct {
	// leveldb handle
	db *leveldb.DB
}

// Open a DB for given path, directory will be created if the path dose
// not exist.
func (db *DB) OpenFile(fileName string, options leveldbOpt.Options) (*DB, error) {
	ldb, err := leveldb.OpenFile(fileName, options)
	if err != nil {
		return nil, err
	}
	db := new(DB)
	db.db = ldb
	return db
}

// Close a DB.
func (db *DB) Close() error {
	return db.db.Close()
}
