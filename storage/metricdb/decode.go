// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"strconv"
)

// Decode db key into metric.
func decodeKey(key []byte, m *models.Metric) error {
	s := string(key)
	if len(s) <= stampLen {
		return ErrCorrupted
	}
	idx := len(s) - stampLen
	str := s[idx:]
	n, err := strconv.ParseUint(str, convBase, 32)
	if err != nil {
		return err
	}
	m.Stamp = horizon + uint32(n)
	return nil
}

// Decode db value into metric.
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
