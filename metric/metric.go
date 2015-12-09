// Copyright 2015 Eleme Inc. All rights reserved.

package metric

// Metric with name and value
type Metric struct {
	Name   string  // metric name
	Stamp  uint64  // metric timestamp
	Value  float64 // metric value
	Score  float64 // metric anomaly score
	AvgOld float64 // previous average value
	AvgNew float64 // current average value
}

// New creates a Metric.
func New() *Metric {
	m := new(Metric)
	m.Stamp = 0
	m.Score = 0
	return m
}
