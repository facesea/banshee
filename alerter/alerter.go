// Copyright 2015 Eleme Inc. All rights reserved.

package alerter

import (
	"encoding/json"
	"os/exec"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/filter"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"github.com/eleme/banshee/util/safemap"
)

// Limit for buffered detected metric results, further results will be dropped
// if this limit is reached.
const bufferedMetricResultsLimit = 10 * 1024

// Alerter alerts on anomalies detected.
type Alerter struct {
	// Storage
	db *storage.DB
	// Config
	cfg *config.Config
	// Filter
	filter *filter.Filter
	// Input
	In chan *models.Metric
	// Alertings stamps
	m *safemap.SafeMap
}

// Alerting message.
type msg struct {
	Project *models.Project `json:"project"`
	Metric  *models.Metric  `json:"metric"`
	User    *models.User    `json:"user"`
}

// New creates a alerter.
func New(cfg *config.Config, db *storage.DB, filter *filter.Filter) *Alerter {
	al := new(Alerter)
	al.cfg = cfg
	al.db = db
	al.filter = filter
	al.In = make(chan *models.Metric, bufferedMetricResultsLimit)
	al.m = safemap.New()
	return al
}

// Start several goroutines to wait for detected metrics, then check each
// metric with all the rules, the configured shell command will be executed
// once a rule is hit.
func (al *Alerter) Start() {
	for i := 0; i < al.cfg.Alerter.Workers; i++ {
		go al.work()
	}
}

// work waits for detected metrics, then check each metric with all the
// rules, the configured shell command will be executed once a rule is hit.
func (al *Alerter) work() {
	for {
		metric := <-al.In
		// Check interval.
		v, ok := al.m.Get(metric.Name)
		if ok && metric.Stamp-v.(uint32) < al.cfg.Alerter.Interval {
			return
		}
		// Test with rules.
		rules := al.filter.MatchedRules(metric)
		for _, rule := range rules {
			// Test
			if !rule.Test(metric) {
				continue
			}
			// Project
			var proj *models.Project
			if err := al.db.Admin.DB().Model(rule).Related(proj); err != nil {
				log.Error("project not found: %v, skiping..", err)
				continue
			}
			// Users
			var users []models.User
			if err := al.db.Admin.DB().Model(proj).Related(&users, "Users"); err != nil {
				log.Error("get users: %v, skiping..", err)
				continue
			}
			// Send
			for _, user := range users {
				d := &msg{
					Project: proj,
					Metric:  metric,
					User:    &user,
				}
				// Exec
				b, _ := json.Marshal(d)
				cmd := exec.Command(al.cfg.Alerter.Command, string(b))
				if err := cmd.Run(); err != nil {
					log.Error("exec %s: %v", al.cfg.Alerter.Command, err)
				}
			}
			if len(users) != 0 {
				al.m.Set(metric.Name, metric.Stamp)
			}
		}
	}
}
