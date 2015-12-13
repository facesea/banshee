// Copyright 2015 Eleme Inc. All rights reserved.

// Anomalies detection algorithm.
//
// The algorithm is based on exponential weighted moving average,
// exponential moving standard deviation and the 3-sigma rule.
//
package algorithm

import "math"

// An exponential weighted moving average is a type of infinite impulse
// response filter that applies weighting factors which decrease exponential.
// The weighting for each older datum decreases exponential, never reaching
// zero.
//
// The EWMA for a series Y can be calculated recursively:
//   S[0] = Y[0]
//   S[i] = Y[i] * F + (1-F) * S[i-1]
// Where the F is the weighting factor, a number between 0 and 1, A higher F
// discounts older observations faster. And S is the moving average.
//
// Function Ewma implements the recurrence formula of this algorithm.
func Ewma(wf, avgOld, v float64) float64 {
	return (1-wf)*avgOld + wf*v
}

// An exponential weighted moving standard deviation is the exactly same with
// ewma, for a series Y it can be calculated recursively:
//  D[0] = 0
//  D[i] = Sqrt((1-F)*D[i-1]*D[i-1] + F*(Y[i]-S[i-1])*(Y[i]-S[i]))
// Where the D is the moving standard deviation and S is the moving average.
//
// Function Ewms implements the recurrence formula of this algorithm.
func Ewms(wf, avgOld, avgNew, stdOld, v float64) float64 {
	return math.Sqrt((1-wf)*stdOld*stdOld + wf*(v-avgOld)*(v-avgNew))
}

// The 3-sigma rule, also named 68-95-99.7 rule tells us that nearlly all
// values (99.7%) lie within 3 standard deviations of the mean in a normal
// distribution. So if a metric dosen't meet this rule, it must be an anomaly.
//
// To describe it in pseudocode:
//   if abs(V-S) > 3*D then
//     return Anomaly
// Where the V is the current value to be detected, the S is the series
// average, and the D is the series standard deviation.
//
// Function Div3Sigma implements this rule and returns the divison result, which
// is named as metric score here.
//
// If the score is larger than 1, the trending is anomalously raising up.
// If the score is smaller than -1, the trending is anomalously reducing
// down.
func Div3Sigma(avg, std, v float64) float64 {
	if std == 0 {
		return 0
	}
	return (v - avg) / (3 * std)
}

// Get the 3-sigma score, which can help to test whether the latest element of
// a series is an anomaly. If the score is larger than 1, the latest element is
// anomalously large, otherwise if it is smaller than -1, it is anomalously
// small.
//
// Note that we don't use this function in our detector, instead, we use the
// recurrence functions and thus the detection complexity is only O(1).
func sigs(wf float64, series []float64) float64 {
	avg := series[0]
	std := 0.0
	for i := 0; i < len(series); i++ {
		v := series[i]
		avgOld := avg
		avg = Ewma(wf, avgOld, v)
		std = Ewms(wf, avgOld, avg, std, v)
	}
	return Div3Sigma(avg, std, series[len(series)-1])
}
