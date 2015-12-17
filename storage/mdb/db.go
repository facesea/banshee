// Copyright 2015 Eleme Inc. All rights reserved.

// Package mdb handles the storage for metrics.
package mdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Prefix
const (
	prefixIndex = 'i'
	prefixStamp = 's'
)

// Timstamp encoding
const (
	// Further incoming timestamp will be stored as the diff to this horizon.
	stampHorizon uint32 = 1450322633
	// Timestamps will be converted to to 36-hex string format before they are
	// put to db.
	stampConvBase = 36
	// A timestamp in 36-hex format string with this length is enough.
	stampStringLength = 7
)

// Options is for db opening.
type Options struct {
	// Metrics expiration.
	Experiation uint32
}

// DB handles the metric storage.
type DB struct {
	// LevelDB
	db *leveldb.DB
	// Expiration
	exp uint32
}

// Open a DB by fileName.
func Open(fileName string, options *Options) (*DB, error) {
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db, options.Experiation}, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// Operations

// Put a metric to db.
func (db *DB) Put(m *models.Metric) error {
	// Put metric and index.
	idxKey := encodeIndexKey(m)
	idxValue := encodeIndexValue(m)
	stampKey := encodeStampKey(m)
	stampValue := encodeStampValue(m)
	batch := new(leveldb.Batch)
	batch.Put(idxKey, idxValue)
	batch.Put(stampKey, stampValue)
	// Remove outdated metrics.
	metrics, err := db.Range(m.Name, stampHorizon+1, m.Stamp-db.exp)
	if err != nil {
		return err
	}
	if len(metrics) != 0 {
		for _, m := range metrics {
			key := encodeStampKey(m)
			batch.Delete(key)
		}
	}
	return db.db.Write(batch, nil)
}

// Get a metric by name and stamp.
func (db *DB) Get(name string, stamp uint32) (*models.Metric, error) {
	m := &models.Metric{Name: name, Stamp: stamp}
	stampKey := encodeStampKey(m)
	value, err := db.db.Get(stampKey, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	err = decodeMetricValue(value, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Range returns metrics in a stamp range, the start and end are both included
// in the range.
func (db *DB) Range(name string, start, end uint32) (metrics []*models.Metric, err error) {
	startM := &models.Metric{Name: name, Stamp: start}
	endM := &models.Metric{Name: name, Stamp: end + 1}
	startStampKey := encodeStampKey(startM)
	endStampKey := encodeStampKey(endM)
	iter := db.db.NewIterator(&util.Range{Start: startStampKey, Limit: endStampKey}, nil)
	for iter.Next() {
		m := &models.Metric{}
		key := iter.Key()
		value := iter.Value()
		// Fill name and stamp
		err = decodeMetricKey(key, m)
		if err != nil {
			return
		}
		// File value, score, and average.
		err = decodeMetricValue(value, m)
		if err != nil {
			return
		}
	}
	err = nil
	return
}

// Indexes returns all metric indexes.
func (db *DB) Indexes() (indexes []*models.MetricIndex, err error) {
	prefix := []byte(fmt.Sprintf("%c", prefixIndex))
	iter := db.db.NewIterator(util.BytesPrefix(prefix), nil)
	for iter.Next() {
		idx := &models.MetricIndex{}
		key := iter.Key()
		value := iter.Value()
		// Fill name
		err = decodeIndexKey(key, idx)
		if err != nil {
			return
		}
		// Fill score, average
		err = decodeIndexValue(value, idx)
		if err != nil {
			return
		}
	}
	err = nil
	return
}
