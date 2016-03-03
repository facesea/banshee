// Copyright 2016 Eleme Inc. All rights reserved.

package mathutil

import (
	"github.com/eleme/banshee/util/assert"
	"math/rand"
	"testing"
)

func TestAverage(t *testing.T) {
	vals := []float64{1, 2, 3, 4}
	assert.Ok(t, Average(vals) == 2.5)
}

func TestStdDev(t *testing.T) {
	vals := []float64{1, 2, 2, 1}
	assert.Ok(t, StdDev(vals, Average(vals)) == .5)
}

func genValues(n int) []float64 {
	var vals []float64
	for i := 0; i < n; i++ {
		vals = append(vals, rand.Float64()*float64(rand.Intn(1024)))
	}
	return vals
}

func BenchmarkAverageNum605(b *testing.B) {
	vals := genValues(605)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Average(vals)
	}
}

func BenchmarkStdDevNum605(b *testing.B) {
	vals := genValues(605)
	avg := Average(vals)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StdDev(vals, avg)
	}
}
