// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import "math"

// Get the average of float64 values.
func average(vals []float64) float64 {
	var sum float64
	for i := 0; i < len(vals); i++ {
		sum += vals[i]
	}
	return sum / float64(len(vals))
}

// Get the standard deviation of float64 values, with
// an input average.
func stdDev(vals []float64, avg float64) float64 {
	var sum float64
	for i := 0; i < len(vals); i++ {
		dis := vals[i] - avg
		sum += dis * dis
	}
	return math.Sqrt(sum / float64(len(vals)))
}
