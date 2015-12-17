package detector

import (
	"sync"

	"github.com/eleme/banshee/models"
)

// cache is two map contain hitCache for Detector with read/write lock
// to keep goroutine safety
type cache struct {
	lockForWLC     *sync.RWMutex
	lockForRHC     *sync.RWMutex
	whiteListCache map[string]bool
	rulesHitCache  map[string]map[string]bool
}

// newCache creates a new hitCache
func newCache() *cache {
	return &cache{
		lockForWLC:     &sync.RWMutex{},
		lockForRHC:     &sync.RWMutex{},
		whiteListCache: make(map[string]bool),
		rulesHitCache:  make(map[string]map[string]bool),
	}

}

// hitWhiteListCache - Check if a metric hit the hitCache--'whiteListCache'
func (c *cache) hitWhiteListCache(m *models.Metric) (hit bool, cache bool) {
	c.lockForWLC.RLock()
	defer c.lockForWLC.RUnlock()
	v, e := c.whiteListCache[m.Name]
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
func (c *cache) setWLC(m *models.Metric, rule *models.Rule, pass bool) {
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.lockForRHC.Lock()
	defer c.lockForRHC.Unlock()

	if rule == nil {
		c.whiteListCache[m.Name] = pass
		return
	}

	_, exists := c.rulesHitCache[rule.Pattern]
	if exists {
		c.rulesHitCache[rule.Pattern][m.Name] = pass
	} else {
		c.rulesHitCache[rule.Pattern] = map[string]bool{m.Name: pass}
	}
	c.whiteListCache[m.Name] = pass
}

// updateRules - update cache by rules
func (c *cache) updateRules(rules []models.Rule) {
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.lockForRHC.Lock()
	defer c.lockForRHC.Unlock()
	delList := []string{}
	for key, value := range c.rulesHitCache {
		needDel := true
		for _, rule := range rules {
			if rule.Pattern == key {
				needDel = false
			}
		}
		if needDel {
			for metric, _ := range value {
				delete(c.whiteListCache, metric)
			}
			delList = append(delList, key)
		}
	}
	for _, v := range delList {
		delete(c.rulesHitCache, v)
	}
}
