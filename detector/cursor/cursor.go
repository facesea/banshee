// Copyright 2015 Eleme Inc. All rights reserved.

// Package cursor implements the algorithm which moves the detection
// state forward.
package cursor

import (
	"math"

	"github.com/eleme/banshee/models"
)

// Cursor is a detection state pusher.
type Cursor struct {
	wf     float64 // Weighted factor
	leastC int     // Least count
}

// New create a cursor.
func New(wf float64, leastC int) *Cursor {
	return &Cursor{wf, leastC}
}

// Next moves the state next with the metric, this will fill in the metric's average
// and score fields and return the next state.
// If the pervious state s is nil, a new state with count 1 and value from
// metric's will be returned as the next state.
func (c *Cursor) Next(s *models.State, m *models.Metric) *models.State {
	n := &models.State{}
	if s == nil {
		// As first
		m.Average = m.Value
		m.Score = 0
		n.Average = m.Value
		n.StdDev = 0
		n.Count = 1
		return n
	}
	// Old average will be kept for refering.
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
	// Don't move forward the stddev if the current metric is anomalous.
	if m.IsAnomalous() {
		wf:=c.wf*s.Average/math.Abs(s.Average-m.Value)
		n.Average = ewma(wf, s.Average, m.Value)
		n.StdDev = s.StdDev
	}
	return n
}
