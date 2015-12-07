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

// Put time series name with a value, this enables using series name to
// value as a hash table.
func (db *DB) PutName(name string, v Value) error {
	key := encodeSeriesNameKey(name)
	val := encodeValue(v)
	return db.db.Put([]byte(key), []byte(val), nil)
}

// Put datapoint into db with name, time stamp and value, additionally
// with the value for series name (default as 0).
func (db *DB) PutWithNameValue(name string, nameValue Value, t TimeStamp, v Value) error {
	err := db.PutName(name, nameValue)
	if err != nil {
		return err
	}
	key := encodeTimeStampKey(name, t)
	val := encodeValue(v)
	err = db.db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		db.DeleteName(name)
		return err
	}
	return nil
}

// Delete time series name from db.
func (db *DB) DeleteName(name string) error {
	key := encodeSeriesNameKey(name)
	err := db.db.Delete([]byte(key), nil)
	return err
}

// Delete datapoint from db by name and timestamp. If this is the last
// datapoint in the series, the series name will be deleted too.
func (db *DB) Delete(name string, t TimeStamp) error {
	// NotImplemented
}

// Get datapoint value from db by name and timestamp.
func (db *DB) Get(name string, t TimeStamp) (Value, error) {
	// NotImplemented
}

// Get the value of series name (as key).
func (db *DB) GetNameValue(name string) (Value, error) {
	// NotImplemented
}
