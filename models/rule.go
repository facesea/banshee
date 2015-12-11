// Copyright 2016 Eleme Inc. All rights reserved.

package models

const (
	WhenTrendUp      = 0x0001
	WhenTrendDown    = 0x0010
	WhenValueAtLeast = 0x0100
	WhenValueAtMost  = 0x1000
)

// Rule is a type to describe alerting rule.
type Rule struct {
	// Pattern is a wildcard string
	Pattern string
	// Condition
	WhenWhat int
	AtLeast  float64
	AtMoset  float64
}
