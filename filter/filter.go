// Copyright 2015 Eleme Inc. All rights reserved.

package filter

import (
	"strings"
	"sync"

	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/safemap"
)

// Filter is a safeMap contain hitCache for Detector
type Filter struct {
	AddRule chan *models.Rule
	DelRule chan *models.Rule
	//children map[string]*childCache
	children *safemap.SafeMap
}

// childFilter is a suffix tree
type childFilter struct {
	lock         *sync.RWMutex
	matchedRules []*models.Rule
	//children map[string]*childCache
	children *safemap.SafeMap
}

// Limit for buffered changed rules
const bufferedChangedRulesLimit = 1000

// NewFilter creates a new filter
func NewFilter() (newFilter *Filter) {
	newFilter = &Filter{
		AddRule:  make(chan *models.Rule, bufferedChangedRulesLimit),
		DelRule:  make(chan *models.Rule, bufferedChangedRulesLimit),
		children: safemap.New(),
	}
	go newFilter.addRules()
	go newFilter.delRules()
	return newFilter
}

// newChildCache creates a new childCache
func newChildFilter() *childFilter {
	return &childFilter{
		lock:         &sync.RWMutex{},
		matchedRules: []*models.Rule{},
		children:     nil,
	}
}

func (c *childFilter) matchedRs(l []string) []*models.Rule {
	if len(l) == 0 {
		c.lock.RLock()
		defer c.lock.RUnlock()
		return c.matchedRules
	}
	rules := []*models.Rule{}
	if c.children == nil {
		return rules
	}
	v, e := c.children.Get("*")
	if e {
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	v, e = c.children.Get(l[0])
	if e {
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	return rules
}

// MatchedRules checks if a metric hit the hitCache, if hit return cache value
func (f *Filter) MatchedRules(m *models.Metric) []*models.Rule {
	rules := []*models.Rule{}
	l := strings.Split(m.Name, ".")
	v, e := f.children.Get("*")
	if e {
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	v, e = f.children.Get(l[0])
	if e {
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	return rules
}

func (f *Filter) addRule(rule *models.Rule) {
	l := strings.Split(rule.Pattern, ".")
	if !f.children.Has(l[0]) {
		f.children.Set(l[0], newChildFilter())
	}
	v, _ := f.children.Get(l[0])
	l = l[1:]
	for len(l) > 0 {
		if v.(*childFilter).children == nil {
			v.(*childFilter).children = safemap.New()
		}
		if v.(*childFilter).children.Has(l[0]) {
			v, _ = v.(*childFilter).children.Get(l[0])
		} else {
			v.(*childFilter).children.Set(l[0], newChildFilter())
			v, _ = v.(*childFilter).children.Get(l[0])
		}
		l = l[1:]
	}
	v.(*childFilter).lock.Lock()
	defer v.(*childFilter).lock.Unlock()
	v.(*childFilter).matchedRules = append(v.(*childFilter).matchedRules, rule)
}

func (f *Filter) delRule(rule *models.Rule) {
	l := strings.Split(rule.Pattern, ".")
	if !f.children.Has(l[0]) {
		return
	}
	v, _ := f.children.Get(l[0])
	l = l[1:]
	for len(l) > 0 {
		if v.(*childFilter).children == nil {
			return
		}
		if v.(*childFilter).children.Has(l[0]) {
			v, _ = v.(*childFilter).children.Get(l[0])
		} else {
			return
		}
		l = l[1:]
	}
	v.(*childFilter).lock.Lock()
	defer v.(*childFilter).lock.Unlock()
	rules := []*models.Rule{}
	for _, r := range v.(*childFilter).matchedRules {
		if !rule.Equal(r) {
			rules = append(rules, r)
		}
	}
	v.(*childFilter).matchedRules = rules
}

// updateRules clean dirty cache
func (f *Filter) addRules() {
	for {
		rule := <-f.AddRule
		f.addRule(rule)
	}
}

func (f *Filter) delRules() {
	for {
		rule := <-f.DelRule
		f.delRule(rule)
	}
}
