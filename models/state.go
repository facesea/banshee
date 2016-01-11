// Copyright 2015 Eleme Inc. All rights reserved.

package models

// State is to record the recurrence result of detection.
type State struct {
	// Detection won't start if the current analyzation suite is too small,
	// Count is to record the count of metrics hit the current grid.
	// It won't be updated until it is greater than startSize.
	// And there is no need to use a lock or sync.Mutex for it since incoming
	// metrics at the same time never duplicates.
	Count uint32
	// Current moving average value for this metric at this grid.
	Average float64
	// Current moving standard deviation for this metric at this grid.
	StdDev float64
}
