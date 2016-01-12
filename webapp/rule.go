// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/eleme/banshee/models"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

// createRule request
type createRuleRequest struct {
	ProjectID             int     `json:"projectID"`
	Pattern               string  `json:"pattern"`
	OnTrendUp             bool    `json:"onTrendUp"`
	OnTrendDown           bool    `json:"onTrendDown"`
	OnValueGt             bool    `json:"onValueGt"`
	OnValueLt             bool    `json:"onValueLt"`
	OnTrendUpAndValueGt   bool    `json:"onTrendUpAndValueGt"`
	OnTrendDownAndValueLt bool    `json:"onTrendDownAndValueLt"`
	ThresholdMax          float64 `json:"thresholdMax"`
	ThresholdMin          float64 `json:"thresholdMin"`
	TrustLine             float64 `json:"trustLine"`
}

// createRule creates a rule.
func createRule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Request
	req := &createRuleRequest{}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Validate
	if len(req.Pattern) <= 0 {
		// Pattern is empty.
		ResponseError(w, ErrRulePattern)
		return
	}
	if req.ProjectID <= 0 {
		// ProjectID is invalid.
		ResponseError(w, ErrRuleProjectID)
		return
	}
	if (req.OnValueGt || req.OnTrendUpAndValueGt) && req.ThresholdMax == 0 {
		// ThresholdMax should not be 0.
		ResponseError(w, ErrRuleThresholdMaxRequired)
		return
	}
	if (req.OnValueLt || req.OnTrendDownAndValueLt) && req.ThresholdMin == 0 {
		// ThresholdMin should not be 0.
		ResponseError(w, ErrRuleThresholdMinRequired)
		return
	}
	// When
	when := 0
	if req.OnTrendUp {
		when = when | models.WhenTrendUp
	}
	if req.OnTrendDown {
		when = when | models.WhenTrendDown
	}
	if req.OnValueGt {
		when = when | models.WhenValueGt
	}
	if req.OnValueLt {
		when = when | models.WhenValueLt
	}
	if req.OnTrendUpAndValueGt {
		when = when | models.WhenTrendUpAndValueGt
	}
	if req.OnTrendDownAndValueLt {
		when = when | models.WhenTrendDownAndValueLt
	}
	if when == 0 {
		// No condition
		ResponseError(w, ErrRuleNoCondition)
		return
	}
	// Find project.
	proj := &models.Project{}
	if err := db.Admin.DB().First(proj, req.ProjectID).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrProjectNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Create rule.
	rule := &models.Rule{
		ProjectID:    req.ProjectID,
		Pattern:      req.Pattern,
		When:         when,
		ThresholdMax: req.ThresholdMax,
		ThresholdMin: req.ThresholdMin,
		TrustLine:    req.TrustLine,
	}
	if err := db.Admin.DB().Create(rule).Error; err != nil {
		// Write errors.
		sqliteErr, ok := err.(sqlite3.Error)
		if ok {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintNotNull:
				ResponseError(w, ErrNotNull)
				return
			case sqlite3.ErrConstraintPrimaryKey:
				ResponseError(w, ErrPrimaryKey)
				return
			case sqlite3.ErrConstraintUnique:
				ResponseError(w, ErrDuplicateRulePattern)
				return
			}
		}
		// Unexcepted error.
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	// Cache
	db.Admin.RulesCache.Put(rule)
	// Response
	ResponseJSONOK(w, rule)
}

// deleteRule deletes a rule from a project.
func deleteRule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		ResponseError(w, ErrProjectID)
		return
	}
	// Delete
	if err := db.Admin.DB().Delete(&models.Rule{ID: id}).Error; err != nil {
		switch err {
		case gorm.RecordNotFound:
			ResponseError(w, ErrRuleNotFound)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Cache
	db.Admin.RulesCache.Delete(id)
}
