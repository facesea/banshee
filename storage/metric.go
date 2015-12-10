// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"github.com/eleme/banshee/errors"
	"github.com/eleme/banshee/metric"
)

// Get a metric from db inplace. example:
//   m := &metric.Metric{Name: "foo", Stamp: 1449740827, Value: 1.43}
//   err := db.GetMetric(m)
//   if err != nil {
//     if errors.IsFatal(err) {
//       log.Fatal("failed to get metric: %v", err)
//     }
//       log.Error("failed to get metric: %v", err)
//   }
//
func (db *DB) GetMetric(m *metric.Metric) error {
	key := db.packMetricKey(m)
	val, err := db.d.Get(key, nil)
	if err != nil {
		if db.IsCorrupted(err) {
			return errors.NewErrFatal(err)
		}
		return err
	}
	return db.unpackMetric(val, m)
}

// Put a metric into db. example:
//   err := db.PutMetric(m)
//   if err != nil {
//     if errors.IsFatal(err) {
//       log.Fatal("failed to put metric: %v", err)
//     }
//       log.Error("failed to put metric: %v", err)
//   }
//
func (db *DB) PutMetric(m *metric.Metric) error {
	key := db.packMetricKey(m)
	val := db.packMetric(m)
	err := db.d.Put(key, val, nil)
	if err != nil {
		if db.IsCorrupted(err) {
			return errors.NewErrFatal(err)
		}
		return err
	}
	return nil
}

// Pack metric key into bytes.
func (db *DB) packMetricKey(m *metric.Metric) []byte {
	numGrid := db.cfg.Periodicity[0]
	grid := db.cfg.Periodicity[1]
	periodicity := numGrid * grid
	gridNo := (m.Stamp % periodicity) / grid
	return []byte(fmt.Sprintf("%s:%d", m.Name, gridNo))
}

// Pack metric db value from metric.
func (db *DB) packMetric(m *metric.Metric) []byte {
	return []byte(fmt.Sprintf("%.3f:%.3f:%d", m.Avg, m.Std, m.Count))
}

// Unpack metric from bytes in place, this will fill in metric average,
// standard deviation and count into the metric.
func (db *DB) unpackMetric(v []byte, m *metric.Metric) error {
	n, err := fmt.Sscanf(string(v), "%f:%f:%d", m.Avg, m.Std, m.Count)
	if err != nil {
		return err
	}
	if n != 3 {
		return err
	}
	return nil
}
