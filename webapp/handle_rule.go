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
	ProjectID    int     `json:"projectID"`
	Pattern      string  `json:"pattern"`
	When         int     `json:"when"`
	ThresholdMax float64 `json:"thresholdMax"`
	ThresholdMin float64 `json:"thresholdMin"`
	TrustLine    float64 `json:"trustLine"`
}

// createRule creates a rule.
func createRule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Request
	req := &createRuleRequest{
		When:      models.WhenTrendUp | models.WhenTrendDown,
		TrustLine: 0,
	}
	if err := RequestBind(r, req); err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Validate
	if len(req.Pattern) <= 0 {
		ResponseError(w, ErrRulePattern)
		return
	}
	if req.When <= 0x1 || req.When >= 0x3f {
		ResponseError(w, ErrRuleWhen)
		return
	}
	if req.ProjectID <= 0 {
		ResponseError(w, ErrRuleProjectID)
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
		When:         req.When,
		ThresholdMax: req.ThresholdMax,
		ThresholdMin: req.ThresholdMin,
		TrustLine:    req.TrustLine,
	}
	if err := db.Admin.DB().Create(rule).Error; err != nil {
		switch err {
		case sqlite3.ErrConstraintNotNull:
			ResponseError(w, ErrNotNull)
			return
		case sqlite3.ErrConstraintPrimaryKey:
			ResponseError(w, ErrPrimaryKey)
			return
		case sqlite3.ErrConstraintUnique:
			ResponseError(w, ErrDuplicateRulePattern)
			return
		default:
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
	}
	// Cache
	db.Admin.RulesCache.Put(rule)
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
	db.Admin.RulesCache.Delete(&models.Rule{ID: id})
}
