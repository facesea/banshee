package detector

import (
	"github.com/eleme/banshee/models"
	"sync"
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
func newCache(rules *[]string) *cache {
	//FIXME use rules
	return &cache{
		lockForWLC:     &sync.RWMutex{},
		lockForRHC:     &sync.RWMutex{},
		whiteListCache: make(map[string]bool),
		rulesHitCache:  make(map[string]map[string]bool),
	}

}

// hitWhiteListCache - Check if a metric hit the hitCache--'whiteListCache'
func (c *cache) hitWhiteListCache(m *models.Metric) (hit bool,cache bool) {
	c.lockForWLC.RLock()
	defer c.lockForWLC.RUnlock()
	v, e := c.whiteListCache[m.Name]
	if e {
		if v {
			return true,true
		}
		return true,false
	}
	return false,false
}

// findRule - Find the rule match metric from Cache when the nil rule param
// send to the other func
func (c *cache) findRule(m *models.Metric) (rulePattern string) {
	//FIXME
	return "nil"
}

// setWLC - Put a white list hitCache into cache , add it to rulesHitCache also
// rule can be nil , when it's nil use O(n) algorithm find rule from rulesHitCache
func (c *cache) setWLC(m *models.Metric, rule *models.Rule,pass bool) {
	if rule == nil {
		//FIXME
	}
	c.lockForWLC.Lock()
	defer c.lockForWLC.Unlock()
	c.lockForRHC.Lock()
	defer c.lockForRHC.Unlock()
	_, exists := c.rulesHitCache[rule.Pattern]
	if exists {
		c.rulesHitCache[rule.Pattern][m.Name] = pass
	} else {
		c.rulesHitCache[rule.Pattern] = map[string]bool{m.Name: pass}
	}
	c.whiteListCache[m.Name] = pass
}


func (c *cache) updateRules(rules []models.Rule) {
	//FIXME
}
