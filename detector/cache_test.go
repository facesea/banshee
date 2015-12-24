// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"testing"

	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
)

func TestCache(t *testing.T) {
	c := newCache()
	ma := &models.Metric{Name: "abcde"}
	mb := &models.Metric{Name: "abcdf"}
	mc := &models.Metric{Name: "abcdg"}
	c.setCache(ma, true)
	c.setCache(mb, false)
	e, v := c.hitCache(ma)
	assert.Ok(t, e && v)
	e, v = c.hitCache(mb)
	assert.Ok(t, e && !v)
	e, v = c.hitCache(mc)
	assert.Ok(t, !e)

}
