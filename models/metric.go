// Copyright 2015 Eleme Inc. All rights reserved.

package models

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

// Is the metric anomalous?
func (m *Metric) IsAnomalous() bool {
	return m.IsAnomalousTrendUp() || m.IsAnomalousTrendDown()
}

// Is the metric trending up?
func (m *Metric) IsTrendUp() bool {
	return m.Score > 0
}

// Is the metric trending down?
func (m *Metric) IsTrendDown() bool {
	return m.Score < 0
}

// Is the metric trending up anomaously?
func (m *Metric) IsAnomalousTrendUp() bool {
	return m.Score > 1
}

// Is the metric trending down anomaously?
func (m *Metric) IsAnomalousTrendDown() bool {
	return m.Score < -1
}
