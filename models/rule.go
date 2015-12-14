// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/util"
)

// Conditions 	000001 = 0x1 	000010 = 0x2	000100 = 0x4
//				001000 = 0x8	010000 = 0x10	100000 = 0x20
const (
	WhenAnomalousTrendUp   = 0x1
	WhenAnomalousTrendDown = 0x2
	//When anomalous greater than thresholdMax AND trend up
	WhenAnomalousTrendUpTo = 0x4
	//When anomalous less than thresholdMin AND trend Down
	WhenAnomalousTrendDownTo = 0x8
	WhenValueGreaterThan     = 0x10
	WhenValueLessThan        = 0x20
)

// Rule is a type to describe alerting rule.
type Rule struct {
	// Pattern is a wildcard string
	Pattern string
	// Condition
	//When = sum (conditions const) , for example :
	//0x15 when WhenAnomalousTrendUp OR WhenAnomalousTrendUpTo OR WhenValueGreaterThan
	When         int
	ThresholdMax float64
	ThresholdMin float64
}

// Test Metric with Rule.
func (rule *Rule) Test(m *Metric) bool {
	if !util.FnMatch(m.Name, rule.Pattern) {
		return false
	}
	test := false
	tmp := rule.When
	for k := 0x20; k >= 0x1; k = k / 2 {
		if tmp >= k {
			tmp = tmp - k
			switch k {
			case WhenAnomalousTrendUp:
				test = test || m.IsAnomalousTrendUp()
			case WhenAnomalousTrendDown:
				test = test || m.IsAnomalousTrendDown()
			case WhenValueGreaterThan:
				test = test || m.Value >= rule.ThresholdMax
			case WhenValueLessThan:
				test = test || m.Value <= rule.ThresholdMin
			case WhenAnomalousTrendUpTo:
				test = test || (m.IsAnomalousTrendUp() && m.Value >= rule.ThresholdMax)
			case WhenAnomalousTrendDownTo:
				test = test || (m.IsAnomalousTrendDown() && m.Value <= rule.ThresholdMin)
			}
		}
		if test {
			return test
		}
	}
	return test
}

// Validate rule.When
func (rule *Rule) IsValid() bool {
	return rule.When > 0x1 && rule.When <= 0x3F
}
