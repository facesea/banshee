// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/util"
	"math/rand"
	"testing"
)

func TestDiv3sigma(t *testing.T) {
	l := make([]float64, 0)
	n := 100
	min := 100.0
	max := 120.0
	rangeN := int(max - min)
	wf := 0.07
	for i := 0; i < n; i++ {
		l = append(l, min+float64(rand.Intn(rangeN))+rand.Float64())
	}
	// The latest element should be normal in this series
	score := sigs(wf, l)
	util.Assert(t, score < 1)
	util.Assert(t, score > -1)
	// The latest element should be a greater one in this series
	l = append(l, 130.0)
	score = sigs(wf, l)
	util.Assert(t, score > 0)
	// The latest element should be a smaller one in this series
	l[n] = 90.0
	score = sigs(wf, l)
	util.Assert(t, score < 0)
	// The latest element should be anomaly one in this series
	l[n] = 140
	score = sigs(wf, l)
	util.Assert(t, score > 1)
	// The latest element should be anomaly one in this series
	// The latest element should be anomaly one in this series
	l[n] = 75
	score = sigs(wf, l)
	util.Assert(t, score < -1)
}
