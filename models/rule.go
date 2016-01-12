// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"fmt"
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/util"
	"path/filepath"
	"strings"
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

// Rule is a type to describe alerter rule.
type Rule struct {
	// Rule may be cached.
	cache `sql:"-" json:"-"`
	// ID in db.
	ID int `gorm:"primary_key" json:"id"`
	// Project belongs to
	ProjectID int `sql:"index;not null" json:"projectID"`
	// Pattern is a wildcard string
	Pattern string `sql:"size:400;not null;unique" json:"pattern"`
	// Alerting condition, described as the logic OR result of basic conditions
	// above:
	//   When = condition1 | condition2
	// example:
	//   0x3 = WhenTrendUp | WhenTrendDown
	When int `json:"when"`
	// Optional thresholds data, only used if the rule condition is about
	// threshold. Although we really don't need any thresholds for trending
	// analyzation and alertings, but we still offer a way to alert by
	// thresholds.
	ThresholdMax float64 `json:"thresholdMax"`
	ThresholdMin float64 `json:"thresholdMin"`
	// Never alert for values below this line. The mainly reason to provide
	// this option is to ignore the volatility of values with a very small
	// average level.
	TrustLine float64 `json:"trustLine"`
	// String representation.
	Repr string `sql:"-" json:"repr"`
}

// IsValid returns true if rule.When is valid.
func (rule *Rule) IsValid() bool {
	rule.RLock()
	defer rule.RUnlock()
	return rule.When >= 0x1 && rule.When <= 0x3F
}

// Copy the rule.
func (rule *Rule) Copy() *Rule {
	dst := &Rule{}
	rule.CopyTo(dst)
	return dst
}

// CopyTo copy the rule to another.
func (rule *Rule) CopyTo(r *Rule) {
	rule.RLock()
	defer rule.RUnlock()
	r.Lock()
	defer r.Unlock()
	r.ID = rule.ID
	r.ProjectID = rule.ProjectID
	r.Pattern = rule.Pattern
	r.When = rule.When
	r.ThresholdMax = rule.ThresholdMax
	r.ThresholdMin = rule.ThresholdMin
	r.TrustLine = rule.TrustLine
}

// Equal tests rule equality
func (rule *Rule) Equal(r *Rule) bool {
	rule.RLock()
	defer rule.RUnlock()
	r.RLock()
	defer rule.RUnlock()
	return (r.ID == rule.ID &&
		r.ProjectID == rule.ProjectID &&
		r.Pattern == rule.Pattern &&
		r.When == rule.When &&
		r.ThresholdMax == rule.ThresholdMax &&
		r.ThresholdMin == rule.ThresholdMin &&
		r.TrustLine == rule.TrustLine)
}

// Test returns true if the metric hits this rule.
func (rule *Rule) Test(m *Metric, cfg *config.Config) bool {
	// RLock if shared.
	rule.RLock()
	defer rule.RUnlock()
	// Ignore it if its value small enough to be trust
	trustLine := rule.TrustLine
	if trustLine == 0 {
		// Check default trustlines
		for pattern, value := range cfg.Detector.DefaultTrustLines {
			if ok, _ := filepath.Match(pattern, m.Name); ok && value != 0 {
				trustLine = value
				break
			}
		}
	}
	// Use rule's trustline
	if m.Value < trustLine {
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

// InitRepr initializes the rule's string repr.
func (rule *Rule) BuildRepr() {
	var parts []string
	if rule.When&WhenTrendUp != 0 {
		parts = append(parts, "trend ↑")
	}
	if rule.When&WhenTrendDown != 0 {
		parts = append(parts, "trend ↓")
	}
	if rule.When&WhenValueGt != 0 {
		s := fmt.Sprintf("value >= %s", util.ToFixed(rule.ThresholdMax, 3))
		parts = append(parts, s)
	}
	if rule.When&WhenValueLt != 0 {
		s := fmt.Sprintf("value <= %s", util.ToFixed(rule.ThresholdMin, 3))
		parts = append(parts, s)
	}
	if rule.When&WhenTrendUpAndValueGt != 0 {
		s := fmt.Sprintf("(trend ↑ && value >= %s)", util.ToFixed(rule.ThresholdMax, 3))
		parts = append(parts, s)
	}
	if rule.When&WhenTrendDownAndValueLt != 0 {
		s := fmt.Sprintf("(trend ↓ && value <= %s)", util.ToFixed(rule.ThresholdMin, 3))
		parts = append(parts, s)
	}
	rule.Repr = strings.Join(parts, " || ")
}
