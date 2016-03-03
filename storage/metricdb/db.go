// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbFilter "github.com/syndtr/goleveldb/leveldb/filter"
	leveldbOpt "github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	// LevelDBBloomFilterBitsPerKey is the leveldb builtin bloom filter
	// bitsPerKey.
	//
	// Filter name will be persisted to disk on a per sstable basis, during
	// reads leveldb will try to find matching filter. And the filter can be
	// replaced after a DB has been created, or say that the filter can be
	// changed between db opens. But if this is done, note that old sstables
	// will continue using the old filter and every new created sstable will
	// use the new filter. For this reason, I make the parameter `bitsPerKey`
	// a constant, there is no need to replace it.
	//
	// Also, this means that no big performance penalty will be experienced
	// when changing the parameter, and the goleveldb docs points this as well.
	//
	// About the probability of false positives, the following link may help:
	// http://pages.cs.wisc.edu/~cao/papers/summary-cache/node8.html
	//
	LevelDBBloomFilterBitsPerKey = 10
)

// Options is db opening options.
type Options struct {
}

// DB handles metrics storage.
type DB struct {
	// LevelDB
	db *leveldb.DB
}

// Open a DB by fileName.
func Open(fileName string) (*DB, error) {
	opts := &leveldbOpt.Options{
		Filter: leveldbFilter.NewBloomFilter(LevelDBBloomFilterBitsPerKey),
	}
	db, err := leveldb.OpenFile(fileName, opts)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// Operations

// Put a metric into db.
func (db *DB) Put(m *models.Metric) error {
	if m.Stamp < horizon {
		return ErrStampTooSmall
	}
	key := encodeKey(m)
	value := encodeValue(m)
	return db.db.Put(key, value, nil)
}

// Get metrics in a timestamp range, the range is left open and right closed.
func (db *DB) Get(name string, start, end uint32) ([]*models.Metric, error) {
	// Key encoding.
	startMetric := &models.Metric{Name: name, Stamp: start}
	endMetric := &models.Metric{Name: name, Stamp: end}
	startKey := encodeKey(startMetric)
	endKey := encodeKey(endMetric)
	// Iterate db.
	iter := db.db.NewIterator(&util.Range{
		Start: startKey,
		Limit: endKey,
	}, nil)
	var metrics []*models.Metric
	for iter.Next() {
		m := &models.Metric{}
		key := iter.Key()
		value := iter.Value()
		// Fill in the name and stamp.
		err := decodeKey(key, m)
		if err != nil {
			return metrics, err
		}
		// Fill in the value, score and average.
		err = decodeValue(value, m)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// Delete metrics in a timestamp range, the range is left open and right
// closed.
func (db *DB) Delete(name string, start, end uint32) (int, error) {
	// Key encoding.
	startMetric := &models.Metric{Name: name, Stamp: start}
	endMetric := &models.Metric{Name: name, Stamp: end}
	startKey := encodeKey(startMetric)
	endKey := encodeKey(endMetric)
	// Iterate db.
	iter := db.db.NewIterator(&util.Range{
		Start: startKey,
		Limit: endKey,
	}, nil)
	batch := new(leveldb.Batch)
	n := 0
	for iter.Next() {
		key := iter.Key()
		batch.Delete(key)
		n++
	}
	if batch.Len() > 0 {
		return n, db.db.Write(batch, nil)
	}
	return n, nil
}

// DeleteTo deletes metrics ranging to a stamp by name.
func (db *DB) DeleteTo(name string, end uint32) (int, error) {
	return db.Delete(name, horizon, end)
}
