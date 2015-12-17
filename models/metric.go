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

// MetricIndex is a container for metric latest name, score and average.
type MetricIndex struct {
	// Name
	Name string
	// Score
	Score float64
	// Average
	Average float64
}

// IsAnomalous test whether the metric is anomalous.
func (m *Metric) IsAnomalous() bool {
	return m.IsAnomalousTrendUp() || m.IsAnomalousTrendDown()
}

// IsTrendUp test whether the metric is trending up.
func (m *Metric) IsTrendUp() bool {
	return m.Score > 0
}

// IsTrendDown test whether the metric is trending down.
func (m *Metric) IsTrendDown() bool {
	return m.Score < 0
}

// IsAnomalousTrendUp test whether the metric is anomalously trending up.
func (m *Metric) IsAnomalousTrendUp() bool {
	return m.Score > 1
}

// IsAnomalousTrendDown test whether the metric is anomalously trending down.
func (m *Metric) IsAnomalousTrendDown() bool {
	return m.Score < -1
}
