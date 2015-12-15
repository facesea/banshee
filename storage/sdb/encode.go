// Copyright 2015 Eleme Inc. All rights reserved.

package sdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
)

func (db *DB) getGridNo(m *models.Metric) int {
	period := db.numGrid * db.gridLen
	return int(m.Stamp%uint32(period)) / db.gridLen
}

// Encode state key by metric.
func (db *DB) encodeKey(m *models.Metric) []byte {
	gridNo := db.getGridNo(m)
	key := fmt.Sprintf("%s:%d", m.Name, gridNo)
	return []byte(key)
}

// Encode state value.
func (db *DB) encodeValue(s *models.State) []byte {
	value := fmt.Sprintf("%.3f:%.3f:%d", s.Average, s.StdDev, s.Count)
	return []byte(value)
}
