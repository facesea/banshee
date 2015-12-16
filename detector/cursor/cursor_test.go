// Copyright 2015 Eleme Inc. All rights reserved.

package cursor

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"math/rand"
	"testing"
)

// Help to generate data suite.
func genMetrics(min, max float64, count int) []*models.Metric {
	arr := make([]*models.Metric, 0)
	delta := max - min
	for i := 0; i < count; i++ {
		value := rand.Float64()*delta + min
		m := &models.Metric{Value: value}
		arr = append(arr, m)
	}
	return arr
}

// Case as first
func TestAsFirst(t *testing.T) {
	wf := 0.05
	leastC := 18
	c := New(wf, leastC)
	m := &models.Metric{Value: 1.32}
	s := c.Next(nil, m)
	assert.Ok(t, m.Average == m.Value)
	assert.Ok(t, m.Score == 0)
	assert.Ok(t, s.Count == 1)
	assert.Ok(t, s.StdDev == 0)
}

// Case count not enough
func TestNotEnough(t *testing.T) {
	wf := 0.05
	leastC := 18
	c := New(wf, leastC)
	s := &models.State{Count: leastC - 1, Average: 0.1, StdDev: 0.1}
	m := &models.Metric{Value: 0.1}
	n := c.Next(s, m)
	assert.Ok(t, m.Score == 0)
	assert.Ok(t, n.Count == leastC)
}

// Simple case.
func TestSimple(t *testing.T) {
	wf := 0.05
	leastC := 18
	c := New(wf, leastC)
	l := genMetrics(120.0, 140.0, leastC)
	var s *models.State = nil
	for _, m := range l {
		s = c.Next(s, m)
		assert.Ok(t, !m.IsAnomalous())
	}
	// Should be normal
	m := &models.Metric{Value: 130.0}
	s = c.Next(s, m)
	assert.Ok(t, !m.IsAnomalous())
	// Should be anomalous up
	m = &models.Metric{Value: 160.0}
	s = c.Next(s, m)
	assert.Ok(t, m.IsAnomalous())
	assert.Ok(t, m.IsAnomalousTrendUp())
	// Should be anomalous down
	m = &models.Metric{Value: 100.0}
	s = c.Next(s, m)
	assert.Ok(t, m.IsAnomalous())
	assert.Ok(t, m.IsAnomalousTrendDown())
}

// Case anomaly after an anomaly.
func TestAnomalyAfterBigAnomaly(t *testing.T) {
	wf := 0.05
	leastC := 18
	c := New(wf, leastC)
	l := genMetrics(120.0, 140.0, 100)
	var s *models.State = nil
	for _, m := range l {
		s = c.Next(s, m)
		assert.Ok(t, !m.IsAnomalous())
	}
	// Give a big anomaly
	m := &models.Metric{Value: 2000}
	s = c.Next(s, m)
	assert.Ok(t, m.IsAnomalousTrendUp())
	// Test another anomaly
	m = &models.Metric{Value: 190}
	s = c.Next(s, m)
	// FIXME: now is trending down
	assert.Ok(t, m.IsAnomalous())
}

// Case slowly trending up.
func TestSlowlyTrendingUp(t *testing.T) {

}
