// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"sync"

	"github.com/eleme/banshee/models"
)

// cache is two map contain hitCache for Detector with read/write lock
// to keep goroutine safety
type cache struct {
	lockForWLC *sync.RWMutex
	cache      map[string]bool
}

// newCache creates a new hitCache
func newCache() *cache {
	return &cache{
		lockForWLC: &sync.RWMutex{},
		cache:      make(map[string]bool),
	}

}

// hitWhiteListCache - Check if a metric hit the hitCache--'whiteListCache'
func (c *cache) hitCache(m *models.Metric) (hit bool, cache bool) {
	c.lockForWLC.RLock()
	defer c.lockForWLC.RUnlock()
	v, e := c.cache[m.Name]
	if e {
		if v {
			return true, true
		}
		return true, false
	}
	return false, false
}

// setWLC - Put a white list hitCache into cache , add it to rulesHitCache also
// rule can be nil , when it's nil the cache should be a blackListHit case , it
// will not be added to rulesHitCache
func (c *cache) setCache(m *models.Metric, pass bool) {
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.cache[m.Name] = pass
}

// updateRules - clean cache
func (c *cache) updateRules() {
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.cache = make(map[string]bool)
}
