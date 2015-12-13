// Copyright 2015 Eleme Inc. All rights reserved.

package models

import "fmt"

// Metric is a data container for time series datapoint.
type Metric struct {
	// Name
	Name string
	// Metric unix time stamp
	Stamp uint32
	// Metric value
	Value float64
	// Anomaly score
	Score float64
	// Average old
	Average float64
}

// Return true if current metric behaves anomalous.
func (m *Metric) IsAnomalous() bool {
	return m.Score >= 1 || m.Score <= -1
}

// Return true if current metric behaves anomalous and the trend is raising up.
func (m *Metric) IsAnomalousTrendUp() bool {
	return m.Score >= 1
}

// Return true if current metric behaves anomalous and the trend is decreasing
// down.
func (m *Metric) IsAnomalousTrendDown() bool {
	return m.Score <= -1
}

// String presentation for a metric.
func (m *Metric) String() string {
	return fmt.Sprintf("name=%s stamp=%d value=%.3f score=%.3f average=%.3f",
		m.Name, m.Stamp, m.Value, m.Score, m.Average)
}
