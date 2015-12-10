// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"fmt"
	"strconv"

	"github.com/eleme/banshee/errors"
	"github.com/eleme/banshee/metric"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	leaderIndex  = 'A' // index table leader char
	leaderMetric = 'B' // metric table leader char
)

var maxIdKey = []byte("MaxID")

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
	id, err := db.getMetricId(m.Name)
	if err != nil {
		return err
	}
	key := db.packMetricKey(id, m.Stamp)
	v, err := db.d.Get(key, nil)
	if err != nil {
		if db.IsCorrupted(err) {
			return errors.NewErrFatal(err)
		}
		return err
	}
	return db.unpackMetric(v, m)
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
	id, err := db.getMetricId(m.Name)
	if err != nil {
		return err
	}
	key := db.packMetricKey(id, m.Stamp)
	err = db.d.Put(key, db.packMetric(m), nil)
	if err != nil {
		if db.IsCorrupted(err) {
			return errors.NewErrFatal(err)
		}
		return err
	}
	return nil
}

// Get metric id by metric name. If not found, new one.
func (db *DB) getMetricId(name string) (uint32, error) {
	key := db.packMetricIdKey(name)
	v, err := db.d.Get(key, nil)
	if err == nil {
		id, err := db.unpackMetricId(v)
		if err != nil {
			// Found
			return 0, err
		}
		return id, nil
	}
	if err != leveldb.ErrNotFound {
		// Unexcepted
		if db.IsCorrupted(err) {
			return 0, errors.NewErrFatal(err)
		}
		return 0, err
	}
	// Not found
	db.dl.Lock()
	defer db.dl.Unlock()
	id := uint32(0)
	v, err = db.d.Get(maxIdKey, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			// First key
			id = 1
		} else {
			// Unexcepted
			if db.IsCorrupted(err) {
				return 0, errors.NewErrFatal(err)
			}
			return 0, err
		}
	} else {
		id, err = db.unpackMetricId(v)
		if err != nil {
			return 0, err
		}
		id += 1
	}
	// Increment and Put
	v = db.packMetricId(id)
	db.d.Put(maxIdKey, v, nil)
	db.d.Put(key, v, nil)
	return id, nil
}

//----------------------------------------
// Pack and unpack
//----------------------------------------

// Pack metric id into bytes.
func (db *DB) packMetricId(id uint32) []byte {
	return []byte(strconv.FormatUint(uint64(id), 36))
}

// Unpack metric id into bytes.
func (db *DB) unpackMetricId(v []byte) (uint32, error) {
	id, err := strconv.ParseUint(string(v), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// Pack metric id key into bytes.
func (db *DB) packMetricIdKey(name string) []byte {
	return []byte(fmt.Sprintf("%c%s", leaderIndex, name))
}

// Pack metric key into bytes.
func (db *DB) packMetricKey(id uint32, stamp uint32) []byte {
	numGrid := db.cfg.Periodicity[0]
	grid := db.cfg.Periodicity[1]
	periodicity := numGrid * grid
	gridNo := (stamp % periodicity) / grid
	return []byte(fmt.Sprintf("%d:%d", id, gridNo))
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
