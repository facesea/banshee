// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/util/assert"
	"math/rand"
	"testing"
)

func generateData() []float64 {
	l := make([]float64, 0)
	n := 100
	min := 100.0
	max := 120.0
	rangeN := int(max - min)
	for i := 0; i < n; i++ {
		l = append(l, min+float64(rand.Intn(rangeN))+rand.Float64())
	}
	return l
}

func TestDiv3sigmaBasic(t *testing.T) {
	wf := 0.07
	l := generateData()
	n := len(l)
	// The latest element should be normal in this series
	score := sigs(wf, l)
	assert.Ok(t, score < 1)
	assert.Ok(t, score > -1)
	// The latest element should be a greater one in this series
	l = append(l, 130.0)
	score = sigs(wf, l)
	assert.Ok(t, score > 0)
	// The latest element should be a smaller one in this series
	l[n] = 90.0
	score = sigs(wf, l)
	assert.Ok(t, score < 0)
	// The latest element should be anomaly one in this series
	l[n] = 140
	score = sigs(wf, l)
	assert.Ok(t, score > 1)
	// The latest element should be anomaly one in this series
	// The latest element should be anomaly one in this series
	l[n] = 75
	score = sigs(wf, l)
	assert.Ok(t, score < -1)
}

func TestDiv3sigmaFactor(t *testing.T) {
	// The smaller the trend factor is, the more sensitive the detection will
	// be.
	l := generateData()
	l = append(l, 130.0)
	l = append(l, 140.0)
	l = append(l, 150.0)
	l = append(l, 160.0)
	assert.Ok(t, sigs(0.03, l) > 1)
	assert.Ok(t, sigs(0.09, l) < 1)
}
