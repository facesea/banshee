// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage/metricdb"
)

// Metrics with longer names will be refused.
const maxMetricNameLen = 256

// Validate input metrics, the following limitations will
// be checked:
//
//	1. Metric name shouldn't be longer than maxMetricNameLen.
//	2. Metric stamp shouldn't be smaller than `horizon`.
//
func validateMetric(m *models.Metric) error {
	if len(m.Name) > maxMetricNameLen {
		// Name too long.
		return ErrMetricNameTooLong
	}
	if m.Stamp < metricdb.Horizon() {
		// Stamp too small.
		return ErrMetricStampTooSmall
	}
	return nil
}
