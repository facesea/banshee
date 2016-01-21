// Copyright 2015 Eleme Inc. All rights reserved.

// Package alerter implements an alerter to send sms/email messages
// on anomalies found.
package alerter

import (
	"encoding/json"
	"os/exec"
	"sync/atomic"
	"time"

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
	// Alertings counters
	c *safemap.SafeMap
}

// Alerting message.
type msg struct {
	Project *models.Project `json:"project"`
	Metric  *models.Metric  `json:"metric"`
	User    *models.User    `json:"user"`
}

// New creates a alerter.
func New(cfg *config.Config, db *storage.DB) *Alerter {
	al := new(Alerter)
	al.cfg = cfg
	al.db = db
	al.In = make(chan *models.Metric, bufferedMetricResultsLimit)
	al.m = safemap.New()
	al.c = safemap.New()
	return al
}

// Start several goroutines to wait for detected metrics, then check each
// metric with all the rules, the configured shell command will be executed
// once a rule is hit.
func (al *Alerter) Start() {
	log.Info("start %d alerter workers..", al.cfg.Alerter.Workers)
	for i := 0; i < al.cfg.Alerter.Workers; i++ {
		go al.work()
	}
	go func() {
		ticker := time.NewTicker(time.Hour * 24)
		for _ = range ticker.C {
			al.c.Clear()
		}
	}()
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
		// Check alert times in one day
		v, ok = al.c.Get(metric.Name)
		if ok && atomic.LoadUint32(v.(*uint32)) > al.cfg.Alerter.OneDayLimit {
			log.Warn("%s hit alerting one day limit, skipping..", metric.Name)
			return
		}
		if !ok {
			var newCounter uint32
			newCounter = 1
			al.c.Set(metric.Name, &newCounter)
		} else {
			atomic.AddUint32(v.(*uint32), 1)
		}
		// Universals
		var univs []models.User
		if err := al.db.Admin.DB().Where("universal = ?", true).Find(&univs).Error; err != nil {
			log.Error("get universal users: %v, skiping..", err)
			continue
		}
		for _, rule := range metric.TestedRules {
			// Project
			proj := &models.Project{}
			if err := al.db.Admin.DB().Model(rule).Related(proj); err != nil {
				log.Error("project, %v, skiping..", err)
				continue
			}
			// Users
			var users []models.User
			if err := al.db.Admin.DB().Model(proj).Related(&users, "Users").Error; err != nil {
				log.Error("get users: %v, skiping..", err)
				continue
			}
			users = append(users, univs...)
			// Send
			for _, user := range users {
				d := &msg{
					Project: proj,
					Metric:  metric,
					User:    &user,
				}
				// Exec
				if len(al.cfg.Alerter.Command) == 0 {
					log.Warn("alert command not configured")
					continue
				}
				b, _ := json.Marshal(d)
				cmd := exec.Command(al.cfg.Alerter.Command, string(b))
				if err := cmd.Run(); err != nil {
					log.Error("exec %s: %v", al.cfg.Alerter.Command, err)
					continue
				}
				log.Info("send message to %s with %s ok", user.Name, metric.Name)
			}
			if len(users) != 0 {
				al.m.Set(metric.Name, metric.Stamp)
			}
		}
	}
}
