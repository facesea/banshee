// Copyright 2016 Eleme Inc. All rights reserved.

// Package mathutil provides math util functions.
package mathutil

import "math"

// Average returns the mean value of float64 values.
func Average(vals []float64) float64 {
	var sum float64
	for i := 0; i < len(vals); i++ {
		sum += vals[i]
	}
	return sum / float64(len(vals))
}

// StdDev returns the standard deviation of float64 values, with an input
// average.
func StdDev(vals []float64, avg float64) float64 {
	var sum float64
	for i := 0; i < len(vals); i++ {
		dis := vals[i] - avg
		sum += dis * dis
	}
	return math.Sqrt(sum / float64(len(vals)))
}
