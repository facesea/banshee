// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"strings"
	"sync"

	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/safemap"
)

// cache is two map contain hitCache for Detector with read/write lock
// to keep goroutine safety
type cache struct {
	wordsLock *sync.RWMutex
	words     map[string]string
	Rc        chan *models.Rule
	//children map[string]*childCache
	children *safemap.SafeMap
}

// childCache is a suffix tree
type childCache struct {
	hit   bool
	cache bool
	//children map[string]*childCache
	children *safemap.SafeMap
}

// Limit for buffered changed rules
const bufferedChangedRulesLimit = 1000

// newCache creates a new hitCache
func newCache() (newcache *cache) {
	newcache = &cache{
		wordsLock: &sync.RWMutex{},
		words:     make(map[string]string),
		Rc:        make(chan *models.Rule, bufferedChangedRulesLimit),
		children:  safemap.New(),
	}
	go newcache.updateRules()
	return newcache
}

// newChildCache creates a new childCache
func newChildCache() *childCache {
	return &childCache{
		hit:      false,
		cache:    false,
		children: nil,
	}
}

// addWord add a new word to word list of the cache
func (c *cache) addWord(word string) string {
	c.wordsLock.Lock()
	defer c.wordsLock.Unlock()
	n := len(c.words)
	s := ""
	for n > 0 {
		k := n % (10 + 26 + 26)
		n = n / (10 + 26 + 26)
		switch {
		case k >= 10+26:
			s = s + string('A'+k-10-26)
		case k >= 10:
			s = s + string('a'+k-10)
		default:
			s = s + string('0'+k)
		}
	}
	c.words[word] = s
	return s
}

// encode the metricName by word list in order to compress the string
func (c *cache) encode(metricName string) []string {
	c.wordsLock.RLock()
	defer c.wordsLock.RUnlock()
	tmp := strings.Split(metricName, ".")
	l := []string{}
	for _, s := range tmp {
		v, e := "*", true
		if s != "*" {
			v, e = c.words[s]
		}
		if !e {
			c.wordsLock.RUnlock()
			v = c.addWord(s)
			c.wordsLock.RLock()
		}
		l = append(l, v)
	}
	return l
}

// hitCache checks if a metric hit the hitCache, if hit return cache value
func (c *cache) hitCache(m *models.Metric) (hit bool, cache bool) {
	l := c.encode(m.Name)
	str := l[0]
	l = l[1:]
	child, e := c.children.Get(str)
	for len(l) > 0 {
		if !e {
			hit, cache = false, false
			return
		}
		str = l[0]
		l = l[1:]
		child, e = child.(*childCache).children.Get(str)
	}
	if !e {
		hit, cache = false, false
		return
	}
	return child.(*childCache).hit, child.(*childCache).cache
}

// setCache put a white list hitCache into cache , add it to rulesHitCache also
// rule can be nil , when it's nil the cache should be a blackListHit case , it
// will not be added to rulesHitCache
func (c *cache) setCache(m *models.Metric, pass bool) {
	l := c.encode(m.Name)
	str := l[0]
	l = l[1:]
	child, e := c.children.Get(str)
	if !e {
		c.children.Set(str, newChildCache())
		child, e = c.children.Get(str)
	}
	for len(l) > 0 {
		if child.(*childCache).children == nil {
			child.(*childCache).children = safemap.New()
		}
		str = l[0]
		l = l[1:]
		if !child.(*childCache).children.Has(str) {
			child.(*childCache).children.Set(str, newChildCache())
		}
		child, _ = child.(*childCache).children.Get(str)
	}
	child.(*childCache).hit = true
	child.(*childCache).cache = pass
}

// cleanChildCache recursive clean the dirty cache of changed rule
func (c *cache) cleanChildCache(l []string, child *childCache) {
	if len(l) == 0 {
		child.hit = false
		child.cache = false
		return
	}
	if child.children == nil {
		return
	}
	str := l[0]
	tmp := l[1:]
	if str == "*" {
		for _, v := range child.children.Items() {
			c.cleanChildCache(tmp, v.(*childCache))
		}
	}
	nextChild, e := child.children.Get(str)
	if e {
		c.cleanChildCache(tmp, nextChild.(*childCache))
	}
}

// cleanCache clean the dirty cache of changed rule
func (c *cache) cleanCache(rule *models.Rule) {
	l := c.encode(rule.Pattern)
	str := l[0]
	l = l[1:]
	if str == "*" {
		for _, v := range c.children.Items() {
			c.cleanChildCache(l, v.(*childCache))
		}
	}
	child, e := c.children.Get(str)
	if e {
		c.cleanChildCache(l, child.(*childCache))
	}
}

// updateRules clean dirty cache
func (c *cache) updateRules() {
	for {
		rule := <-c.Rc
		c.cleanCache(rule)
	}
}
