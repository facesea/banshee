package tsdb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type TimeStamp uint64 // timestamp in seconds
type Value float64    // value in float64

type DB struct {
	db *leveldb.DB
}

func OpenFile(fileName string, options leveldb.Options) (*DB, error) {
	db, err := leveldb.OpenFile(fileName, options)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) Put(name string, t TimeStamp, v Value) error {
	ns := EncodeSeriesNameKey(name)
	err := db.db.Put(ns, []byte{}, nil)
	if err != nil {
		return err
	}
	ts := EncodeTimeStampKey(name, t)
	err := db.db.Put(ts, []byte(EncodeValue(v)), nil)
	if err != nil {
		db.db.Delete(ns, nil)
		return err
	}
	return nil
}
