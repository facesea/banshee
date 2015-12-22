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

// Alerter will goroutine to alert metric which match rule
type Alerter struct {
	db       *storage.DB
	cfg      *config.Config
	rc       chan *models.Metric
	stampMap *safemap.SafeMap
}

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
	alerter.rc = make(chan *models.Metric, bufferedMetricResultsLimit)
	alerter.stampMap = safemap.New()
	return alerter
}

// Start - start several goroutines to wait for detected metrics,
// then check each metric with all the rules, the configured shell command will
// be executed once a rule is hit.
func (alerter *Alerter) Start() {
	for i := 0; i < alerter.cfg.Alerter.Workers; i++ {
		go alerter.work()
	}
}

// work - wait for detected metrics, then check each metric with all the
// rules, the configured shell command will be executed once a rule is hit.
func (alerter *Alerter) work() {
	for {
		metric := <-alerter.rc

		v, e := alerter.stampMap.Get(metric.Name)
		if e && metric.Stamp-uint32(v) < alerter.cfg.Interval {
			return
		}
		alerter.stampMap.Set(metric.Name, metric.Stamp)

		rules := alerter.db.Admin.Rules()
		for _, rule := range rules {
			if rule.Test(metric) {
				proj, err := alerter.db.Admin.GetProject(rule.ProjectID)
				if err != nil {
					log.Error("%v, projectID: %d ruleid: %d", err, rule.ID, rule.ProjectID)
					continue
				}

				for _, user := range proj.Users {
					msg := &msg{
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
