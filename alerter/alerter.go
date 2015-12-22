// Copyright 2015 Eleme Inc. All rights reserved.

package alerter

import (
	"encoding/json"
	"os/exec"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
)

// Limit for buffered detected metric results, further results will be dropped
// if this limit is reached.
const bufferedMetricResultsLimit = 10 * 1024

// Alerter will goroutine to alert metric which match rule
type Alerter struct {
	db  *storage.DB
	cfg *config.Config
	rc  chan *models.Metric
}

type alertMsg struct {
	Project *models.Project `json:"project"`
	Metric  *models.Metric  `json:"metric"`
	User    *models.User    `json:"user"`
}

// New creates a alerter.
func New(cfg *config.Config, db *storage.DB) *Alerter {
	alerter := new(Alerter)
	alerter.cfg = cfg
	alerter.db = db
	alerter.rc = make(chan *models.Metric, bufferedMetricResultsLimit)
	return alerter
}

// StartAlertingWorkers - start several goroutines to wait for detected metrics,
// then check each metric with all the rules, the configured shell command will
// be executed once a rule is hit.
func (alerter *Alerter) StartAlertingWorkers() {
	for i := 0; i < alerter.cfg.Alerter.Workers; i++ {
		go alerter.alertingWork()
	}
}

// alertingWork - wait for detected metrics, then check each metric with all the
// rules, the configured shell command will be executed once a rule is hit.
func (alerter *Alerter) alertingWork() {
	metric := <-alerter.rc
	rules := alerter.db.Admin.Rules()
	for _, rule := range rules {
		if rule.Test(metric) {
			proj, err := alerter.db.Admin.GetProject(rule.ProjectID)
			if err != nil {
				log.Error("getProject by id failed, projectID: %d ruleid: %d", rule.ID, rule.ProjectID)
			} else {
				for _, user := range proj.Users {
					msg := &alertMsg{
						Project: proj,
						Metric:  metric,
						User:    user}
					msgBytes, _ := json.Marshal(msg)
					err := exec.Command(alerter.cfg.Alerter.Command, string(msgBytes)).Run()
					if err != nil {
						log.Error("exec alert command failed : %s", err)
					}
				}
			}
		}
	}
}

// Alert add metric to the alerting chan
func (alerter *Alerter) Alert(metric *models.Metric) {
	select {
	case alerter.rc <- metric:
	default:
		log.Warn("buffered metric results channel is full, drop current metric..")
	}
}
