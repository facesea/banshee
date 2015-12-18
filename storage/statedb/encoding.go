// Copyright 2015 Eleme Inc. All rights reserved.

package statedb

import (
	"fmt"
	"github.com/eleme/banshee/models"
)

// getGirdNo returns the grid number for the metric.
func (db *DB) getGridNo(m *models.Metric) int {
	// GirdNo = (Stamp % Period) / GirdLen
	period := db.numGrid * db.gridLen
	return int(m.Stamp%uint32(period)) / db.gridLen
}

// encodeKey encodes state key by metric.
func (db *DB) encodeKey(m *models.Metric) []byte {
	// Key format is Name:GirdNo
	gridNo := db.getGridNo(m)
	key := fmt.Sprintf("%s:%d", m.Name, gridNo)
	return []byte(key)
}

// encodeValue encodes state value.
func (db *DB) encodeValue(s *models.State) []byte {
	// Value format is Average:StdDev:Count
	value := fmt.Sprintf("%.5f:%.5f:%d", s.Average, s.StdDev, s.Count)
	return []byte(value)
}

// decodeValue decodes db value into state.
func (db *DB) decodeValue(value []byte) (*models.State, error) {
	s := &models.State{}
	n, err := fmt.Sscanf(string(value), "%f:%f:%d", &s.Average, &s.StdDev, &s.Count)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, ErrCorrupted
	}
	return s, nil
}
