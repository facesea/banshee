// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/config"
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

func TestRuleTest(t *testing.T) {
	m := &Metric{Name: "foo", Score: 1.2, Value: 3.567}
	// WhenTrendUp
	r1 := &Rule{When: WhenTrendUp}
	assert.Ok(t, r1.Test(m, nil))
	// WhenTrendDown
	r2 := &Rule{When: WhenTrendDown}
	assert.Ok(t, !r2.Test(m, nil))
	// WhenValueGt
	r3 := &Rule{When: WhenValueGt, ThresholdMax: 1.2}
	assert.Ok(t, r3.Test(m, nil))
	// WhenValueLt
	r4 := &Rule{When: WhenValueLt, ThresholdMin: 1.2}
	assert.Ok(t, !r4.Test(m, nil))
	// WhenTrendUpAndValueGt
	r5 := &Rule{When: WhenTrendUpAndValueGt, ThresholdMax: 1.2}
	assert.Ok(t, r5.Test(m, nil))
	r6 := &Rule{When: WhenTrendUpAndValueGt, ThresholdMax: 9.0}
	assert.Ok(t, !r6.Test(m, nil))
	// WhenTrendUp | WhenTrendDownAndValueLt
	r7 := &Rule{When: WhenTrendUp | WhenTrendDownAndValueLt, ThresholdMin: 2.0}
	assert.Ok(t, r7.Test(m, nil))
	// TrustLine
	r8 := &Rule{When: WhenTrendUp, TrustLine: 1.0}
	assert.Ok(t, r8.Test(m, nil))
	r9 := &Rule{When: WhenTrendUp, TrustLine: 8.0}
	assert.Ok(t, !r9.Test(m, nil))
	// Default TrustLines
	cfg := config.New()
	cfg.Detector.DefaultTrustLines["fo*"] = 4.0
	r10 := &Rule{When: WhenTrendUp}
	assert.Ok(t, !r10.Test(m, cfg))
}
