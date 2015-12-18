// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"strconv"
)

// Stamp
const (
	// Futher timestamps will be stored as the diff to this horizon for storage
	// cost reason.
	horizon uint32 = 1450322633
	// Timestamps will be converted to 36-hex string format before they are put
	// to db, also for storage cost reason.
	convBase = 36
	// A 36-hex string format timestamp with this length is enough to use for
	// later 100 years.
	stampLen = 7
)

// encodeKey encodes db key from metric.
func encodeKey(m *models.Metric) []byte {
	// Key format is Name+Stamp.
	t := m.Stamp - horizon
	v := strconv.FormatUint(uint64(t), convBase)
	s := m.Name + v
	return []byte(s)
}

// encodeValue encodes db value from metric.
func encodeValue(m *models.Metric) []byte {
	// Value format is Value:Score:Average.
	s := fmt.Sprintf("%.5f:%.5f:%.5f", m.Value, m.Score, m.Average)
	return []byte(s)
}
