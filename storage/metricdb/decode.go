// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"strconv"
)

// decodeKey decodes db key into metric, this will fill metric name and metric
// stamp.
func decodeKey(key []byte, m *models.Metric) error {
	s := string(key)
	if len(s) <= stampLen {
		return ErrCorrupted
	}
	// First substring is Name.
	idx := len(s) - stampLen
	m.Name = s[:idx]
	// Last substring is Stamp.
	str := s[idx:]
	n, err := strconv.ParseUint(str, convBase, 32)
	if err != nil {
		return err
	}
	m.Stamp = horizon + uint32(n)
	return nil
}

// decodeValue decodes db value into metric, this will fill metric value,
// average and stddev.
func decodeValue(value []byte, m *models.Metric) error {
	n, err := fmt.Sscanf(string(value), "%f:%f:%f", &m.Value, &m.Score, &m.Average)
	if err != nil {
		return err
	}
	if n != 3 {
		return ErrCorrupted
	}
	return nil
}
