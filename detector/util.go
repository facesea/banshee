// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import "math"

// getMean returns the mean of values.
func getMean(values []float64) float64 {
	sum := float64(0)
	for i := 0; i < len(values); i++ {
		sum += values[i]
	}
	return sum / float64(len(values))
}

// getStd returns the stddev of values.
func getStd(values []float64, mean float64) float64 {
	sum := float64(0)
	for i := 0; i < len(values); i++ {
		dis := values[i] - mean
		sum += dis * dis
	}
	return math.Sqrt(sum / float64(len(values)))
}
