package tsdb

import (
	"strconv"
)

// Increment timeseries length by n, if the n is negative, this
// function will perform decrementing.
func (db *DB) incrLen(name string, n int) error {
	db.l.Lock()
	defer db.l.Unlock()
	key := encodeTsName(name)
	val, err := db.db.Get(key)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	var v int
	if err != leveldb.ErrNotFound {
		v, err = strconv.Atoi(string(val))
		if err != nil {
			return err
		}
		v += n
	} else {
		v = n
	}
	s := strconv.Itoa(v)
	return db.db.Put(key, []byte(s))
}

// Put datapoint into db with name, timestamp and value.
func (db *DB) Put(name string, t TimeStamp, v Value) error {
	key := encodeTsKey(name, t)
	val := encodeTsValue(v)
	err := db.db.Put([]byte(key), []byte(val), nil)
	if err != nil {
		return err
	}
	return db.incrLen(name, 1)
}

// Delete datapoint from db with name, timestamp.
func (db *DB) Delete(name string, t TimeStamp) error {
	key := encodeTsKey(name, t)
	err := db.db.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return db.incrLen(name, -1)
}
