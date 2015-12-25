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
	db.load()
	return db, nil
}

// Close the DB.
func (db *DB) Close() error {
	return db.db.Close()
}

// load indexes from db to cache.
func (db *DB) load() {
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

// get index by name.
func (db *DB) get(name string) (*models.Index, bool) {
	v, ok := db.m.Get(name)
	if ok {
		return v.(*models.Index), true
	}
	return nil, false
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
	// Use an copy.
	idx = idx.Copy()
	// Add to cache.
	idx.Share()
	db.m.Set(idx.Name, idx)
	return nil
}

// Get an index by index name.
func (db *DB) Get(idx *models.Index) error {
	i, ok := db.get(idx.Name)
	if !ok {
		return ErrNotFound
	}
	i.CopyTo(idx)
	return nil
}

// Delete an index by name.
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
	for k, v := range m {
		idx := v.(*models.Index)
		name := k.(string)
		if util.Match(pattern, name) {
			l = append(l, idx.Copy())
		}
	}
	return l
}

// All returns all indexes.
func (db *DB) All() (l []*models.Index) {
	m := db.m.Items()
	for _, v := range m {
		idx := v.(*models.Index)
		l = append(l, idx.Copy())
	}
	return l
}
