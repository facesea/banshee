// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage/metricdb"
)

// Metrics with too long name will be refused.
const MaxMetricNameLen = 256

// validateMetric validates input metric.
func validateMetric(m *models.Metric) error {
	if len(m.Name) > MaxMetricNameLen {
		// Name too long.
		return ErrMetricNameTooLong
	}
	if m.Stamp < metricdb.Horizon() {
		// Stamp too small.
		return ErrMetricStampTooSmall
	}
	return nil
}
