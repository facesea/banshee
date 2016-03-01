// Copyright 2016 Eleme Inc. All rights reserved.

package alerter

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestHourInRange(t *testing.T) {
	assert.Ok(t, hourInRange(3, 0, 6))
	assert.Ok(t, !hourInRange(7, 0, 6))
	assert.Ok(t, !hourInRange(6, 0, 6))
	assert.Ok(t, hourInRange(23, 19, 10))
	assert.Ok(t, hourInRange(6, 19, 10))
	assert.Ok(t, !hourInRange(13, 19, 10))
}

func TestMakeEventID(t *testing.T) {
	// Metric with the same name but different stamps.
	m1 := &models.Metric{Name: "foo", Stamp: 1456815973}
	m2 := &models.Metric{Name: "foo", Stamp: 1456815974}
	assert.Ok(t, makeEventID(m1) != makeEventID(m2))
	// Metric with the same stamp but different names.
	m1 = &models.Metric{Name: "foo", Stamp: 1456815973}
	m2 = &models.Metric{Name: "bar", Stamp: 1456815973}
	assert.Ok(t, makeEventID(m1) != makeEventID(m2))
}
