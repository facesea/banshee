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
	// Pattern is a wildcard string
	Pattern string
	// Condition
	// When = condition1 | condition2, example:
	//   0x3 = WhenTrendUp | WhenTrendDown
	When int
	// Additional
	ThresholdMax float64
	ThresholdMin float64
	// Never alert for values below this line.
	TrustLine float64
}

// Return true if the metric is against this rule.
func (rule *Rule) Test(m *Metric) bool {
	if !util.Match(m.Name, rule.Pattern) {
		return false
	}
	// Ignore it if it's value small enough to be trust
	if m.Value < rule.TrustLine {
		return false
	}

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

// Validate rule.When
func (rule *Rule) IsValid() bool {
	return rule.When >= 0x1 && rule.When <= 0x3F
}
