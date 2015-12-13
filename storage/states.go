// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/syndtr/goleveldb/leveldb"
)

// A statesDB handles the leveldb for detection states storage.
type statesDB struct {
	// LevelDB
	db *leveldb.DB
	// Periodicity
	gridLen  int
	numGrids int
}

// Open a statesDB for detection states storage. If the path dosen't exist, it
// will be created.
func openStatesDB(fileName string, numGrids, gridLen int) (*statesDB, error) {
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}
	return &statesDB{
		db:       db,
		gridLen:  gridLen,
		numGrids: numGrids,
	}, nil
}

// Close a statesDB.
func (db *statesDB) close() error {
	return db.db.Close()
}

// Get current gridNo.
func (db *statesDB) getGridNo(m *models.Metric) int {
	periodicity := db.gridLen * db.numGrids
	return int(m.Stamp%uint32(periodicity)) / db.gridLen
}

// Get the db key of a state by metric.
func (db *statesDB) keyOf(m *models.Metric) []byte {
	gridNo := db.getGridNo(m)
	return []byte(fmt.Sprintf("%s:%d", m.Name, gridNo))
}

// Get the db value of a state.
func (db *statesDB) valOf(s *models.State) []byte {
	return []byte(fmt.Sprintf("%.3f:%.3f:%d", s.Average, s.StdDev, s.Count))
}

// Parse db value bytes into state object.
func (db *statesDB) parse(v []byte) (*models.State, error) {
	s := &models.State{}
	n, err := fmt.Sscanf(string(v), "%f:%f:%d", &s.Average, &s.StdDev, &s.Count)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, ErrCorrupted
	}
	return s, nil
}

// Get the detection state for the metric at the current grid.
//   m := &models.Metric{Name: "foo", Stamp: .., Value: ..}
//   state, err := db.GetState(m)
//   if err != nil {..}
//
func (db *DB) GetState(m *models.Metric) (*models.State, error) {
	key := db.s.keyOf(m)
	val, err := db.s.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return db.s.parse(val)
}

// Put a detection state for a metric to db.
//   m := &models.Metric{Name: "foo", Stamp: .., Value: ..}
//   s := &models.State{Average: .., StdDev: .., Count: ..}
//   err := db.PutState(m)
//   if err != nil {..}
//
func (db *DB) PutState(m *models.Metric, s *models.State) error {
	key := db.s.keyOf(m)
	val := db.s.valOf(s)
	err := db.s.db.Put(key, val, nil)
	if err != nil {
		return err
	}
	return nil
}
