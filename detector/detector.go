// Copyright 2015 Eleme Inc. All rights reserved.

// Package detector implements a tcp server that detects whether incoming
// metrics are anomalies and send alertings on anomalies found.
package detector

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/eleme/banshee/alerter"
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/detector/cursor"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/storage/statedb"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/log"
)

// Detector is a tcp server to detect anomalies.
type Detector struct {
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
	// Alerter
	alerter *alerter.Alerter
	// Filter
	hitCache *cache
	// Cursor
	cursor *cursor.Cursor
}

// New creates a detector.
func New(cfg *config.Config, db *storage.DB, alerter *alerter.Alerter) *Detector {
	d := new(Detector)
	d.cfg = cfg
	d.db = db
	d.alerter = alerter
	d.hitCache = newCache()
	d.db.Admin.SetRuleChan(d.hitCache.Rc)
	d.cursor = cursor.New(cfg.Detector.Factor, cfg.LeastC())
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
			d.alerter.Alert(m)
		}
	}
}

// Test whether a metric matches the rules.
func (d *Detector) match(m *models.Metric) bool {
	// Check rules.
	//FIXME use channel
	rules := d.db.Admin.Rules()
	d.hitCache.updateRules()

	// Check cache first.
	hit, v := d.hitCache.hitCache(m)
	if hit {
		return v
	}

	// Check blacklist.
	for _, pattern := range d.cfg.Detector.BlackList {
		if util.Match(pattern, m.Name) {
			d.hitCache.setCache(m, false)
			log.Debug("%s hit black pattern %s", m.Name, pattern)
			return false
		}
	}

	for _, rule := range rules {
		if util.Match(rule.Pattern, m.Name) {
			d.hitCache.setCache(m, true)
			return true
		}
	}
	// No rules was hit.
	log.Debug("%s hit no rules", m.Name)
	d.hitCache.setCache(m, false)
	return false
}

// Detect incoming metric with 3-sigma rule and fill the metric.Score.
func (d *Detector) detect(m *models.Metric) error {
	// Get pervious state.
	s, err := d.db.State.Get(m)
	if err != nil && err != statedb.ErrNotFound {
		return err
	}
	// Move state next.
	var n *models.State
	if err == statedb.ErrNotFound {
		n = d.cursor.Next(nil, m)
	} else {
		n = d.cursor.Next(s, m)
	}
	// Put the next state to db.
	return d.db.State.Put(m, n)
}
