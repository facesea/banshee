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
	c.setWLC(ma, true)
	c.setWLC(mb, false)
	e, v := c.hitWhiteListCache(ma)
	assert.Ok(t, e && v)
	e, v = c.hitWhiteListCache(mb)
	assert.Ok(t, e && !v)
	e, v = c.hitWhiteListCache(mc)
	assert.Ok(t, !e)

	c.updateRules()
	e, v = c.hitWhiteListCache(ma)
	assert.Ok(t, !e)
	e, v = c.hitWhiteListCache(mb)
	assert.Ok(t, !e)
	c.setWLC(mc, true)
	e, v = c.hitWhiteListCache(mc)
	assert.Ok(t, e && v)

}
