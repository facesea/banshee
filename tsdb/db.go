package tsdb

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

// Time stamps are uint64 numbers (in seconds)
type TimeStamp uint64

// Series values are float64 numbers
type Value float64

// DB is a tsdb database
type DB struct {
	// Leveldb handle
	db *leveldb.DB
}

// Default value for series name key.
const DefaultNameValue Value = 0

// Open or create a DB for the given path.
func OpenFile(fileName string, options leveldb.Options) (*DB, error) {
	db, err := leveldb.OpenFile(fileName, options)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}
}

// Close a DB.
func (db *DB) Close() {
	db.db.Close()
}

// Put datapoint into db with name, time stamp and value. This will
// put 2 key into leveldb: name to 0 and time stamp to value.
func (db *DB) Put(name string, t TimeStamp, v Value) error {
	return PutWithNameValue(name, DefaultNameValue, t, v)
}

// Put datapoint into db with name, time stamp and value, additionally
// with the value for series name (default as 0).
func (db *DB) PutWithNameValue(name string, nameValue Value, t TimeStamp, v Value) error {
	ns := encodeSeriesNameKey(name)
	err := db.db.Put(ns, []byte(encodeValue(nameValue)), nil)
	if err != nil {
		return err
	}
	ts := encodeTimeStampKey(name, t)
	err := db.db.Put(ts, []byte(encodeValue(v)), nil)
	if err != nil {
		db.db.Delete(ns, nil)
		return err
	}
	return nil
}
