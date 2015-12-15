package detector

import (
	"github.com/eleme/banshee/models"
	"sync"
)

// Cache is two map contain hitCache for Detector with read/write lock
// to keep goroutine safety
type cache struct {
	lockForWLC     *sync.RWMutex
	lockForRHC     *sync.RWMutex
	whiteListCache map[string]bool
	rulesHitCache  map[string]map[string]bool
}

// NewCache creates a new hitCache
func newCache(rules []string) *cache {
	//FIXME
	return &cache{
		lockForWLC:     &sync.RWMutex{},
		lockForRHC:     &sync.RWMutex{},
		whiteListCache: make(map[string]bool),
		rulesHitCache:  make(map[string]map[string]bool),
	}

}

// Check if a metric hit the hitCache--'whiteListCache'
func (c *cache) hitWhiteListCache(m models.Metric) bool {
	c.lockForWLC.RLock()
	defer c.lockForWLC.RUnlock()
	v, e := c.whiteListCache[m.Name]
	if e && v {
		return true
	}
	return false
}

// Find the rule match metric from Cache when the nil rule param
// send to the other func
func (c *cache) findRule(m *models.Metric) (rulePattern string) {
	return nil
}

// Put a white list hitCache into cache , add it to rulesHitCache also
// rule can be nil , when it's nil use O(n) algorithm find rule from rulesHitCache
func (c *cache) putWLC(m *models.Metric, rule *models.Rule) {
	if rule == nil {
		//FIXME

	}
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.lockForRHC.Lock()
	defer c.lockForRHC.Unlock()
	_, exists := c.rulesHitCache[rule.Pattern]
	if exists {
		c.rulesHitCache[rule.Pattern][m.Name] = true
	} else {
		c.rulesHitCache[rule.Pattern] = map[string]bool{m.Name: true}
	}
	c.whiteListCache[m.Name] = true
}

// Delete a white list hitCache from cache , del it from rulesHitCache also
// rule can be nil , when it's nil use O(n) algorithm find rule from rulesHitCache
func (c *cache) delWLC(m *models.Metric, rule *models.Rule) {
	if rule == nil {
		//FIXME
	}
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.lockForRHC.Lock()
	defer c.lockForRHC.Unlock()
	delete(c.whiteListCache, m.Name)
	_, exists := c.rulesHitCache[rule.Pattern]
	if exists {
		if _, e := c.rulesHitCache[rule.Pattern][m.Name]; e {
			delete(c.rulesHitCache[rule.Pattern], m.Name)
		}
	}
	delete(c.whiteListCache, m.Name)
}

func (c *cache) updateRules(rules []models.Rule) {

}
