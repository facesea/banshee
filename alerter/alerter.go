// Copyright 2015 Eleme Inc. All rights reserved.

package alerter

import (
	"encoding/json"
	"os/exec"

	"github.com/eleme/banshee/config"
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
func New(cfg *config.Config, db *storage.DB) *Alerter {
	alerter := new(Alerter)
	alerter.cfg = cfg
	alerter.db = db
	alerter.In = make(chan *models.Metric, bufferedMetricResultsLimit)
	alerter.m = safemap.New()
	return alerter
}

// Start several goroutines to wait for detected metrics, then check each
// metric with all the rules, the configured shell command will be executed
// once a rule is hit.
func (alerter *Alerter) Start() {
	for i := 0; i < alerter.cfg.Alerter.Workers; i++ {
		go alerter.work()
	}
}

// work waits for detected metrics, then check each metric with all the
// rules, the configured shell command will be executed once a rule is hit.
func (alerter *Alerter) work() {
	for {
		metric := <-alerter.In
		// Check interval.
		v, ok := alerter.m.Get(metric.Name)
		if ok && metric.Stamp-v.(uint32) < alerter.cfg.Alerter.Interval {
			return
		}
		// Test with rules.
		var projs []*models.Project
		alerter.db.Admin.GetProjects(&projs)
		for _, proj := range projs {
			for _, rule := range proj.Rules {
				if rule.Test(metric) {
					// Tested ok.
					for _, user := range proj.Users {
						// Send message to each user.
						d := &msg{
							Project: proj,
							Metric:  metric,
							User:    user,
						}
						b, _ := json.Marshal(d)
						cmd := exec.Command(alerter.cfg.Alerter.Command, string(b))
						if err := cmd.Run(); err != nil {
							log.Error("exec %s: %v", alerter.cfg.Alerter.Command, err)
						}
					}
					// Add to stamp map.
					alerter.m.Set(metric.Name, metric.Stamp)
				}
			}
		}
	}
}
