// Copyright 2015 Eleme Inc. All rights reserved.

// Package filter implements fast wildcard like filtering based on suffix
// tree.
package filter

import (
	"strings"
	"sync"

	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"github.com/eleme/banshee/util/safemap"
)

// Filter is to filter metrics by rules.
type Filter struct {
	// Rule changes
	addRuleCh chan *models.Rule
	delRuleCh chan *models.Rule
	// Children
	children *safemap.SafeMap
}

// childFilter is a suffix tree.
type childFilter struct {
	lock         *sync.RWMutex
	matchedRules []*models.Rule
	children     *safemap.SafeMap
}

// Limit for buffered changed rules
const bufferedChangedRulesLimit = 1000

// New creates a filter.
func New() *Filter {
	return &Filter{
		addRuleCh: make(chan *models.Rule, bufferedChangedRulesLimit),
		delRuleCh: make(chan *models.Rule, bufferedChangedRulesLimit),
		children:  safemap.New(),
	}
}

// Init from db.
func (f *Filter) Init(db *storage.DB) {
	log.Debug("init rules from cache..")
	// Listen rules changes.
	db.Admin.RulesCache.OnAdd(f.addRuleCh)
	db.Admin.RulesCache.OnDel(f.delRuleCh)
	go f.addRules()
	go f.delRules()
	// Add rules from cache
	var rules []*models.Rule
	db.Admin.RulesCache.All(&rules)
	for _, rule := range rules {
		f.addRule(rule)
	}
}

// newChildCache creates a new childCache
func newChildFilter() *childFilter {
	return &childFilter{
		lock:         &sync.RWMutex{},
		matchedRules: []*models.Rule{},
		children:     nil,
	}
}

// MatchedRules checks if a metric hit, l is the unchecked words list of the metric in order
func (c *childFilter) matchedRs(l []string) []*models.Rule {
	// when len(l)==0 means all words are checked and passed, return all matched rules
	if len(l) == 0 {
		c.lock.RLock()
		defer c.lock.RUnlock()
		return c.matchedRules
	}
	rules := []*models.Rule{}
	//when next level is nil,return empty rules slice
	if c.children == nil {
		return rules
	}
	//check if this level has a "*" node
	v, exist := c.children.Get("*")
	if exist {
		//when has a "*" node, the suffix tree matched the metric words by now, so goto next
		// level and append matched rules to slice
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	//check if this level has a same word node
	v, exist = c.children.Get(l[0])
	if exist {
		//when has the node, matched by now, goto next level and append matched rules to slice
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	//no matched node return empty rules slice, else return all matched rules
	return rules
}

// MatchedRules checks if a metric hit the hitCache, if hit return all hit rules
func (f *Filter) MatchedRules(m *models.Metric) []*models.Rule {
	//split the metric into ordered words
	rules := []*models.Rule{}
	l := strings.Split(m.Name, ".")
	//check if root of the rules suffix tree has a "*" node
	v, exist := f.children.Get("*")
	if exist {
		//when root has a "*" node, goto next level
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	//check if root of the rules suffix tree has the same node to the first word of the metric
	v, exist = f.children.Get(l[0])
	if exist {
		//when has the same word node, goto next level
		rules = append(rules, v.(*childFilter).matchedRs(l[1:])...)
	}
	return rules
}

// addRule add a rule to the suffix tree
func (f *Filter) addRule(rule *models.Rule) {
	//split the rule.Pattern into ordered words
	l := strings.Split(rule.Pattern, ".")
	//if suffix tree root do not has the same node, add it
	if !f.children.Has(l[0]) {
		f.children.Set(l[0], newChildFilter())
	}
	//check if suffix has the same word of the pattern by level step, if not add it
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

// delRule delete a rule from the suffix tree
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

// addRules waits and add new rule to filter.
func (f *Filter) addRules() {
	for {
		rule := <-f.addRuleCh
		f.addRule(rule)
	}
}

// delRules waits and delete rule from filter.
func (f *Filter) delRules() {
	for {
		rule := <-f.delRuleCh
		f.delRule(rule)
	}
}
