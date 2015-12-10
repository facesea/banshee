// Copyright 2015 Eleme Inc. All rights reserved.

package metric

// Metric with name and value
type Metric struct {
	// Name
	Name string
	// Timestamp in seconds, able to use for 90 years from now.
	Stamp uint32
	// Current value
	Value float64 // metric value
	// Current anomaly score
	Score float64
	// Current standard deviation
	Std float64
	// Current average
	Avg float64
	// Previous average
	AvgOld float64
	// Current datapoints count
	Count uint32
}

// New creates a Metric.
func New() *Metric {
	m := new(Metric)
	m.Stamp = 0
	m.Score = 0
	return m
}
