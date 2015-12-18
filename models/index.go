// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Index is for metric indexing, it records metric name for faster metric
// indexing. And metric latest informations like timstamp, score, average are
// also indexed.
type Index struct {
	// Metric name
	Name string
	// Latest stamp for the metric.
	Stamp uint32
	// Latest score for the metric.
	Score float64
	// Latest average for the metric.
	Average float64
}
