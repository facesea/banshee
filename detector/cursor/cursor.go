// Copyright 2015 Eleme Inc. All rights reserved.

// Package cursor implements the algorithm which moves the detection
// state forward.
package cursor

import (
	"math"

	"github.com/eleme/banshee/models"
)

// Cursor is a detection state walker.
type Cursor struct {
	// Weight factor for ewma and ewms, for all metrics. The larger wf is, the
	// faster ewma and ewms moves.
	wf float64
	// The least count for a metric to start the detection, 3-sigma rule will
	// not start to work until the metric's count is large enough.
	leastC int
}

// New create a cursor.
func New(wf float64, leastC int) *Cursor {
	return &Cursor{wf, leastC}
}

// Next moves the state forward for the metric, this operation fills in the
// metric's score and average, and returns the next state.
//
// If the pervious state is nil, the function regards that there's no state if
// for the metric, it returns an initialization state.
//
// If the metric count in its series is not large enough to start the
// detection, the metric will be trusted by setting its score to zero directly.
//
// If the metric behaves normal, the new state will absorb both stddev and
// average from the metric, how much depends on wf. If the metric behaves as an
// anomaly, the new state will not absorb stddev from the metric, and absorb less
// average than usual. Which means that cursor always follows the trending of
// the metric but rejects drastic trend changes.
//
func (c *Cursor) Next(s *models.State, m *models.Metric) *models.State {
	n := &models.State{}
	if s == nil {
		// No previous state, initialize one.
		m.Average = m.Value
		m.Score = 0
		n.Average = m.Value
		n.StdDev = 0
		n.Count = 1
		return n
	}
	// Old average will be saved for later comparison.
	m.Average = s.Average
	// Move forward via ewma & ewms.
	n.Average = ewma(c.wf, s.Average, m.Value)
	n.StdDev = ewms(c.wf, s.Average, n.Average, s.StdDev, m.Value)
	// Fill the score
	if s.Count < c.leastC {
		// Count not enough, trust it
		n.Count = s.Count + 1
		m.Score = 0
	} else {
		// Count enough, get its score
		n.Count = s.Count
		m.Score = div3Sigma(n.Average, n.StdDev, m.Value)
	}
	if m.IsAnomalous() {
		// Absorb average from the usual less than usual.
		wf := c.wf * s.Average / math.Abs(s.Average-m.Value)
		n.Average = ewma(wf, s.Average, m.Value)
		// Don't move forward the stddev, use previous.
		n.StdDev = s.StdDev
	}
	return n
}
