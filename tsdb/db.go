package tsdb

import (
	"sync"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// DB is a tsdb database.
type DB struct {
	// LevelDB
	lv *leveldb.DB
	// Close
	closed uint32
}

// Open a DB for the given path, the DB will be created if not exist. Note
// that the DB can be opened for only once.
//
// OpenFile will return an error with type of ErrCorrupted if corruption
// detected in the DB.
//
// The DB instance is goroutine-safe.
func OpenFile(fileName string, options *opt.Options) (*DB, error) {
	lv, err := leveldb.OpenFile(fileName, options)
	if err != nil {
		return nil, NewErrCorrupted(err)
	}
	db := new(DB)
	db.lv = db
	db.closed = 0
	return db, nil
}

// Close a DB.
func (db *DB) Close() error {
	if db.isClosed() {
		return ErrClosed()
	}
	return db.db.Close()
}

// Check whether DB was closed.
func (db *DB) isClosed() bool {
	return atomic.LoadUint32(&db.closed) != 0
}

// Check whether DB is ok.
func (db *DB) ok() error {
	if db.isClosed() {
		return ErrClosed
	}
	return nil
}

// Batch creates a write batch, which enables write multiple operations once.
func (db *DB) Batch() *Batch {
	return NewBatch(db)
}
