// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"strconv"
)

// Function decodeStampKey decodes db key into metric.
func decodeMetricKey(key []byte, m *models.Metric) error {
	s := string(key)
	if len(s) <= stampStringLength {
		return ErrCorrupted
	}
	idx := len(s) - stampStringLength
	m.Name = s[1:idx]
	stampString := s[idx:]
	n, err := strconv.ParseUint(stampString, stampConvBase, 32)
	if err != nil {
		return err
	}
	m.Stamp = stampHorizon + uint32(n)
	return nil
}

// Function decodeStampValue decodes db value into metric.
func decodeMetricValue(value []byte, m *models.Metric) error {
	n, err := fmt.Sscanf(string(value), "%f:%f:%f", &m.Value, &m.Score, &m.Average)
	if err != nil {
		return err
	}
	if n != 3 {
		return ErrCorrupted
	}
	return nil
}

// Function decodeIndexKey decodes db key into index.
func decodeIndexKey(key []byte, idx *models.MetricIndex) error {
	s := string(key)
	if len(s) < 1 {
		return ErrCorrupted
	}
	idx.Name = s[1:]
	return nil
}

// Function decodeIndexValue decodes db value into index.
func decodeIndexValue(value []byte, idx *models.MetricIndex) error {
	n, err := fmt.Sscanf(string(value), "%f:%f", &idx.Score, &idx.Average)
	if err != nil {
		return err
	}
	if n != 2 {
		return ErrCorrupted
	}
	return nil
}
