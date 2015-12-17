// Copyright 2015 Eleme Inc. All rights reserved.

package mdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"strconv"
)

// Function encodeStampKey encodes stamp key for the metric.
func encodeStampKey(m *models.Metric) []byte {
	t := m.Stamp - stampHorizon
	v := strconv.FormatUint(uint64(t), stampConvBase)
	s := fmt.Sprintf("%c%s%s", prefixStamp, m.Name, v)
	return []byte(s)
}

// Function encodeIndexKey encodes index key for the metric.
func encodeIndexKey(m *models.Metric) []byte {
	s := fmt.Sprintf("%c%s", prefixIndex, m.Name)
	return []byte(s)
}

// Function encodeStampValue encodes stamp value for the metric.
func encodeStampValue(m *models.Metric) []byte {
	s := fmt.Sprintf("%.5f:%.5f:%.5f", m.Value, m.Score, m.Average)
	return []byte(s)
}

// Function encodeNameValue encodes index value for the metric.
func encodeIndexValue(m *models.Metric) []byte {
	s := fmt.Sprintf("%.5f:%.5f", m.Score, m.Average)
	return []byte(s)
}
