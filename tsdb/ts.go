package tsdb

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrTsLen = NewErrCorruptedWithString("invalid time series length")
)

// Increment timeseries length by n, if the n is negative, this
// function will perform decrementing.
func (db *DB) incrLen(name string, n int) error {
	db.tsLenLock.Lock()
	defer db.tsLenLock.Unlock()
	key := encodeTsName(name)
	val, err := db.db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return NewErrCorrupted(err)
	}
	var v int
	if err != leveldb.ErrNotFound {
		v, err = strconv.Atoi(string(val))
		if err != nil {
			return NewErrCorrupted(err)
		}
		v += n
	} else {
		v = n
	}
	if n < 0 {
		return ErrTsLen
	}
	s := strconv.Itoa(v)
	err = db.db.Put([]byte(key), []byte(s), nil)
	if err != nil {
		return NewErrCorrupted(err)
	}
	return err
}

// Put datapoint into db with name, timestamp and value.
func (db *DB) Put(name string, t uint64, v float64) error {
	key := encodeTsKey(name, t)
	val := encodeTsValue(v)
	err := db.db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		return NewErrCorrupted(err)
	}
	return db.incrLen(name, 1)
}

// Delete datapoint from db with name, timestamp. If the key was not found in
// db, ErrNotFound will be returned.
func (db *DB) Delete(name string, t uint64) error {
	key := encodeTsKey(name, t)
	err := db.db.Delete([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return ErrNotFound
		}
		return NewErrCorrupted(err)
	}
	return db.incrLen(name, -1)
}
