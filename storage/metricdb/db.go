// Copyright 2015 Eleme Inc. All rights reserved.

// Package metricdb handles the storage for metrics.
package metricdb

import (
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// DB handles metrics storage.
type DB struct {
	// LevelDB
	db *leveldb.DB
}

// Open a DB by fileName.
func Open(fileName string) (*DB, error) {
	db, err := leveldb.OpenFile(fileName, nil)
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
	key := encodeKey(m)
	value := encodeValue(m)
	return db.db.Put(key, value, nil)
}

// Get metrics in a timstamp range, the range is left open and right closed.
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
func (db *DB) Delete(name string, start, end uint32) error {
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
	for iter.Next() {
		key := iter.Key()
		batch.Delete(key)
	}
	return db.db.Write(batch, nil)
}
