// Copyright 2015 Eleme Inc. All rights reserved.

// Detector is a tcp server which detects whether incoming metrics are
// anomalies and send alertings on anomalies found.
package detector

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/storage/sdb"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/log"
	"github.com/eleme/banshee/util/safemap"
)

// Limit for buffered detected metric results, further results will be dropped
// if this limit is reached.
const bufferedMetricResultsLimit = 10 * 1024

// Detector is a tcp server to detect anomalies.
type Detector struct {
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
	// Results
	rc chan *models.Metric
	// Filter
	matched *safemap.SafeMap
}

// Create a detector.
func New(cfg *config.Config, db *storage.DB) *Detector {
	d := new(Detector)
	d.cfg = cfg
	d.db = db
	d.rc = make(chan *models.Metric, bufferedMetricResultsLimit)
	d.matched = safemap.New()
	return d
}

// Start detector.
func (d *Detector) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", d.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to bind tcp://%s: %v", addr, err)
	}
	log.Info("listening on tcp://%s..", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("failed to accept new conn: %v", err)
		}
		go d.handle(conn)
	}
}

// Handle a connection, it will filter the mertics by rules and detect whether
// the metrics are anomalies.
func (d *Detector) handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	defer func() {
		conn.Close()
		log.Info("conn %s disconnected", addr)
	}()
	log.Info("conn %s established", addr)
	// Scan line by line.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Info("read conn: %v, closing it..", err)
			break
		}
		startAt := time.Now()
		line := scanner.Text()
		m, err := parseMetric(line)
		if err != nil {
			if len(line) > 10 {
				line = line[:10]
			}
			log.Error("parse '%s': %v, skipping..", line, err)
			continue
		}
		if d.match(m) {
			err = d.detect(m)
			if err != nil {
				log.Error("detect: %v, skipping..", err)
				continue
			}
			elapsed := time.Since(startAt)
			log.Debug("name=%s average=%.3f score=%.3f cost=%dμs", m.Name, m.Average, m.Score, elapsed.Nanoseconds()/1000)
			d.rc <- m
		}
	}
}

// Test whether a metric matches the rules.
func (d *Detector) match(m *models.Metric) bool {
	// Check cache first.
	v, ok := d.matched.Get(m.Name)
	if ok {
		return v.(bool)
	}
	// Check blacklist.
	for _, pattern := range d.cfg.Detector.BlackList {
		if util.Match(m.Name, pattern) {
			d.matched.Set(m.Name, false)
			log.Debug("%s hit black pattern %s", m.Name, pattern)
			return false
		}
	}
	// Check rules.
	rules := d.db.UsingA().GetRules()
	for _, rule := range rules {
		if util.Match(m.Name, rule.Pattern) {
			d.matched.Set(m.Name, true)
			return true
		}
	}
	// No rules was hit.
	log.Debug("%s hit no rules", m.Name)
	return false
}

// Detect incoming metric with 3-sigma rule and fill the metric.Score.
func (d *Detector) detect(m *models.Metric) error {
	// Arguments
	wf := d.cfg.Detector.Factor
	startSize := d.cfg.Detector.StartSize
	// Get pervious state.
	s, err := d.db.UsingS().Get(m)
	if err != nil && err != sdb.ErrNotFound {
		return err
	}
	// Next state.
	next := &models.State{}
	if err == sdb.ErrNotFound {
		// Not found, initialize as first
		m.Average = m.Value
		next.Average = m.Value
		next.StdDev = 0
		next.Count = 1
	} else {
		// Found, move to next
		m.Average = s.Average
		next.Average = ewma(wf, s.Average, m.Value)
		next.StdDev = ewms(wf, s.Average, next.Average, s.StdDev, m.Value)
		if s.Count < startSize {
			next.Count = s.Count + 1
		} else {
			next.Count = s.Count
		}
	}
	// Don't calculate the score if current count is not enough.
	if next.Count >= startSize {
		m.Score = div3Sigma(next.Average, next.StdDev, m.Value)
	} else {
		m.Score = 0
	}
	// Don't move forward the standard deviation if the current metric is
	// anomalous.
	if m.IsAnomalous() {
		next.StdDev = s.StdDev
	}
	// Put the next state to db.
	return d.db.UsingS().Put(m, next)
}
