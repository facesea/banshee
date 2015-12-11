// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/util"
)

// Conditions
const (
	WhenAnomalousTrend                  = 0x0
	WhenAnomalousTrendUp                = 0x1
	WhenAnomalousTrendDown              = 0x2
	WhenAnomalousTrendUpTo              = 0x3
	WhenAnomalousTrendDownTo            = 0x4
	WhenAnomalousTrendUpToOrTrendDown   = 0x5
	WhenAnomalousTrendDownToOrTrendUp   = 0x6
	WhenAnomalousTrendUpToOrTrendDownTo = 0x7
	WhenValueGreaterThan                = 0x8
	WhenValueLessThan                   = 0x9
	WhenValueGreaterThanOrLessThan      = 0xa
)

// Rule is a type to describe alerting rule.
type Rule struct {
	// Pattern is a wildcard string
	Pattern string
	// Condition
	When         int
	ThresholdMax float64
	ThresholdMin float64
}

// Test Metric with Rule.
func (rule *Rule) Test(m *Metric) bool {
	if !util.FnMatch(m.Name, rule.Pattern) {
		return false
	}
	switch rule.When {
	case WhenAnomalousTrend:
		return m.IsAnomalous()
	case WhenAnomalousTrendUp:
		return m.IsAnomalousTrendUp()
	case WhenAnomalousTrendDown:
		return m.IsAnomalousTrendDown()
	case WhenAnomalousTrendUpTo:
		return m.IsAnomalousTrendUp() && m.Value >= rule.ThresholdMin
	case WhenAnomalousTrendDownTo:
		return m.IsAnomalousTrendDown() && m.Value <= rule.ThresholdMax
	case WhenAnomalousTrendUpToOrTrendDown:
		return (m.IsAnomalousTrendUp() && m.Value >= rule.ThresholdMin) || m.IsAnomalousTrendDown()
	case WhenAnomalousTrendDownToOrTrendUp:
		return (m.IsAnomalousTrendDown() && m.Value <= rule.ThresholdMax) || m.IsAnomalousTrendUp()
	case WhenAnomalousTrendUpToOrTrendDownTo:
		return (m.IsAnomalousTrendUp() && m.Value >= rule.ThresholdMin) || (m.IsAnomalousTrendDown() && m.Value <= rule.ThresholdMax)
	case WhenValueGreaterThan:
		return m.Value >= rule.ThresholdMin
	case WhenValueLessThan:
		return m.Value <= rule.ThresholdMax
	case WhenValueGreaterThanOrLessThan:
		return m.Value <= rule.ThresholdMax || m.Value >= rule.ThresholdMin
	}
	return false
}
