// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"github.com/eleme/banshee/models"
	"github.com/julienschmidt/httprouter"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type indexByScore []*models.Index

func (l indexByScore) Len() int { return len(l) }

func (l indexByScore) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func (l indexByScore) Less(i, j int) bool {
	now := time.Now().Unix()
	// by `score / ((now - stamp + 2) ^ 1.5)`
	return l[i].Score/math.Pow(float64(uint32(2+now)-l[i].Stamp), 1.5) <
		l[j].Score/math.Pow(float64(uint32(2+now)-l[j].Stamp), 1.5)
}

// getMetricIndexes returns metric names.
func getMetricIndexes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Options
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 50
	}
	order := r.URL.Query().Get("sort")
	if order != "up" && order != "down" {
		order = "up"
	}
	projID, err := strconv.Atoi(r.URL.Query().Get("project"))
	if err != nil {
		projID = 0
	}
	pattern := r.URL.Query().Get("pattern")
	// Index
	var idxs []*models.Index
	if projID > 0 {
		// Rules
		var rules []models.Rule
		if err := db.Admin.DB().Model(&models.Project{ID: projID}).Related(&rules).Error; err != nil {
			ResponseError(w, NewUnexceptedWebError(err))
			return
		}
		// Filter
		for _, rule := range rules {
			idxs = append(idxs, db.Index.Filter(rule.Pattern)...)
		}
	} else {
		// Filter
		idxs = db.Index.Filter(pattern)
	}
	// Sort
	sort.Sort(indexByScore(idxs))
	if order == "down" {
		for i := 0; 2*i < len(idxs); i++ {
			idxs[len(idxs)-1-i], idxs[i] = idxs[i], idxs[len(idxs)-1-i]
		}
	}
	// Limit
	ResponseJSON(w, idxs[:limit])
}

// getMetrics returns metric values.
func getMetrics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Params
	name := ps.ByName("name")
	start, err := strconv.ParseUint(ps.ByName("start"), 10, 32)
	if err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	stop, err := strconv.ParseUint(ps.ByName("stop"), 10, 32)
	if err != nil {
		ResponseError(w, ErrBadRequest)
		return
	}
	// Query
	metrics, err := db.Metric.Get(name, uint32(start), uint32(stop))
	if err != nil {
		ResponseError(w, NewUnexceptedWebError(err))
		return
	}
	ResponseJSON(w, metrics)
}
