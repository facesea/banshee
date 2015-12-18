// Copyright 2015 Eleme Inc. All rights reserved.

// Package indexdb handles the storage for indexes.
package indexdb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/safemap"
	"github.com/syndtr/goleveldb/leveldb"
)

// DB handles indexes storage.
type DB struct {
	// LevelDB
	db *leveldb.DB
	// Cache
	m *safemap.SafeMap
}

// Open a DB by fileName.
func Open(fileName string) (*DB, error) {
	ldb, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}
	m := safemap.New()
	db := new(DB)
	db.db = ldb
	db.m = m
	db.loadM()
	return db, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// loadM loads indexes from db to cache.
func (db *DB) loadM() {
	// Scan values to memory.
	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		// Decode
		key := iter.Key()
		value := iter.Value()
		idx := &models.Index{}
		idx.Name = string(key)
		err := decode(value, idx)
		if err != nil {
			// Skip corrupted values
			continue
		}
		db.m.Set(idx.Name, idx)
	}
}

// Operations.

// Put an index into db.
func (db *DB) Put(idx *models.Index) error {
	// Save to db.
	key := []byte(idx.Name)
	value := encode(idx)
	err := db.db.Put(key, value, nil)
	if err != nil {
		return err
	}
	// Add to cache.
	db.m.Set(idx.Name, idx)
	return nil
}

// Get an index by name.
func (db *DB) Get(name string) (*models.Index, error) {
	v, ok := db.m.Get(name)
	if ok {
		// Found in cache.
		return v.(*models.Index), nil
	}
	// Not found in cache, query db.
	value, err := db.db.Get([]byte(name), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			// Not found
			return nil, ErrNotFound
		}
		// Unexcepted error
		return nil, err
	}
	// Decode
	idx := &models.Index{Name: name}
	err = decode(value, idx)
	if err != nil {
		return nil, err
	}
	// Add to cache.
	db.m.Set(name, idx)
	return idx, nil
}

// Delete an index.
func (db *DB) Delete(name string) error {
	if db.m.Has(name) {
		// Delete in cache.
		db.m.Delete(name)
	}
	key := []byte(name)
	return db.db.Delete(key, nil)
}

// Filter indexes by pattern.
func (db *DB) Filter(pattern string) (l []*models.Index) {
	m := db.m.Items()
	for _, v := range m {
		idx := v.(*models.Index)
		if util.Match(pattern, idx.Name) {
			l = append(l, idx)
		}
	}
	return l
}
