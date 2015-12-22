// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/util"
)

// Conditions
const (
	WhenTrendUp             = 0x1  // 0b000001
	WhenTrendDown           = 0x2  // 0b000010
	WhenValueGt             = 0x4  // 0b000100
	WhenValueLt             = 0x8  // 0b001000
	WhenTrendUpAndValueGt   = 0x10 // 0b010000
	WhenTrendDownAndValueLt = 0x20 // 0b100000
)

// Rule is a type to describe alerting rule.
type Rule struct {
	// Rule may be cached.
	cache `sql:"-"`
	// ID in db.
	ID int `gorm:"primary_key"`
	// Project belongs to
	ProjectID int `sql:"index;not null"`
	// Pattern is a wildcard string
	Pattern string `sql:"size:400;not null;unique"`
	// Alerting condition, described as the logic OR result of basic conditions
	// above:
	//   When = condition1 | condition2
	// example:
	//   0x3 = WhenTrendUp | WhenTrendDown
	When int
	// Optional thresholds data, only used if the rule condition is about
	// threshold. Although we really don't need any thresholds for trending
	// analyzation and alertings, but we still offer a way to alert by
	// thresholds.
	ThresholdMax float64
	ThresholdMin float64
	// Never alert for values below this line. The mainly reason to provide
	// this option is to ignore the volatility of values with a very small
	// average level.
	TrustLine float64
}

// IsValid returns true if rule.When is valid.
func (rule *Rule) IsValid() bool {
	return rule.When >= 0x1 && rule.When <= 0x3F
}

// Copy the rule.
func (rule *Rule) Copy() *Rule {
	if rule.IsShared() {
		rule.RLock()
		defer rule.RUnlock()
	}
	return &Rule{
		ID:           rule.ID,
		ProjectID:    rule.ProjectID,
		Pattern:      rule.Pattern,
		When:         rule.When,
		ThresholdMax: rule.ThresholdMax,
		ThresholdMin: rule.ThresholdMin,
		TrustLine:    rule.TrustLine,
	}
}

// Test returns true if the metric hits this rule.
func (rule *Rule) Test(m *Metric) bool {
	if !util.Match(m.Name, rule.Pattern) {
		// Not match this rule.
		return false
	}

	// Ignore it if it's value small enough to be trust
	if m.Value < rule.TrustLine {
		return false
	}

	// Match conditions.
	ok := false
	if !ok && (rule.When&WhenTrendUp != 0) {
		ok = m.IsAnomalousTrendUp()
	}
	if !ok && (rule.When&WhenTrendDown != 0) {
		ok = m.IsAnomalousTrendDown()
	}
	if !ok && (rule.When&WhenValueGt != 0) {
		ok = m.Value >= rule.ThresholdMax
	}
	if !ok && (rule.When&WhenValueLt != 0) {
		ok = m.Value <= rule.ThresholdMin
	}
	if !ok && (rule.When&WhenTrendUpAndValueGt != 0) {
		ok = m.IsAnomalousTrendUp() && m.Value >= rule.ThresholdMax
	}
	if !ok && (rule.When&WhenTrendDownAndValueLt != 0) {
		ok = m.IsAnomalousTrendDown() && m.Value <= rule.ThresholdMin
	}
	return ok
}
