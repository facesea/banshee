// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestRuleBuildRepr(t *testing.T) {
	var rule *Rule
	// WhenTrendUp
	rule = &Rule{When: WhenTrendUp}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "trend ↑")
	// WhenTrendDown
	rule = &Rule{When: WhenTrendDown}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "trend ↓")
	// WhenTrendUp | WhenTrendDown
	rule = &Rule{When: WhenTrendUp | WhenTrendDown}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "trend ↑ || trend ↓")
	// WhenValueGt
	rule = &Rule{
		When:         WhenValueGt,
		ThresholdMax: 3.1478,
	}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "value >= 3.148")
	// WhenTrendUpAndValueGt
	rule = &Rule{
		When:         WhenTrendUpAndValueGt,
		ThresholdMax: 1.29,
	}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "(trend ↑ && value >= 1.29)")
	// WhenTrendDown | WhenTrendUpAndValueGt
	rule = &Rule{
		When:         WhenTrendUpAndValueGt | WhenTrendDown,
		ThresholdMax: 2223.8,
	}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "trend ↓ || (trend ↑ && value >= 2223.8)")
	// WhenTrendUpAndValueGt | WhenTrendDownAndValueLt
	rule = &Rule{
		When:         WhenTrendUpAndValueGt | WhenTrendDownAndValueLt,
		ThresholdMax: 18987,
		ThresholdMin: 781,
	}
	rule.BuildRepr()
	assert.Ok(t, rule.Repr == "(trend ↑ && value >= 18987) || (trend ↓ && value <= 781)")
}
