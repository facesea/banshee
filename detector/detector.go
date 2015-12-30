// Copyright 2015 Eleme Inc. All rights reserved.

// Package detector implements a tcp server that detects whether incoming
// metrics are anomalies and send alertings on anomalies found.
package detector

import (
	"bufio"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/detector/cursor"
	"github.com/eleme/banshee/filter"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/storage/statedb"
	"github.com/eleme/banshee/util/log"
)

// Detector is a tcp server to detect anomalies.
type Detector struct {
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
	// Filter
	filter *filter.Filter
	// Cursor
	cursor *cursor.Cursor
	// Output
	outs []chan *models.Metric
}

// New creates a detector.
func New(cfg *config.Config, db *storage.DB, filter *filter.Filter) *Detector {
	d := new(Detector)
	d.cfg = cfg
	d.db = db
	d.filter = filter
	d.cursor = cursor.New(cfg.Detector.Factor, cfg.LeastC())
	return d
}

// Out outputs detected metrics to given channel.
func (d *Detector) Out(ch chan *models.Metric) {
	d.outs = append(d.outs, ch)
}

// output detected metrics to outs.
func (d *Detector) output(m *models.Metric) {
	for _, ch := range d.outs {
		select {
		case ch <- m:
		default:
			log.Error("output channel is full, skipping..")
		}
	}
}

// Start detector.
func (d *Detector) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", d.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to bind tcp://%s: %v", addr, err)
	}
	log.Info("listen on tcp://%s..", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept conn: %v", err)
		}
		go d.handle(conn)
	}
}

// Handle a connection, it will filter the mertics by rules and detect whether
// the metrics are anomalies.
func (d *Detector) handle(conn net.Conn) {
	// New conn
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
		// Parse
		line := scanner.Text()
		m, err := parseMetric(line)
		if err != nil {
			if len(line) > 10 {
				line = line[:10]
			}
			log.Error("parse '%s': %v, skipping..", line, err)
			continue
		}
		// Filter
		if d.match(m) {
			// Detect
			err = d.detect(m)
			if err != nil {
				log.Error("failed to detect: %v, skipping..", err)
				continue
			}
			elapsed := time.Since(startAt)
			log.Debug("%dμs %s %.3f", elapsed.Nanoseconds()/1000, m.Name, m.Score)
			// Output
			d.output(m)
			// Store
			if err := d.store(m); err != nil {
				log.Error("store metric %s: %v, skiping..", m.Name, err)
			}
		}
	}
}

// Test whether a metric matches the rules.
func (d *Detector) match(m *models.Metric) bool {
	// Check rules.
	rules := d.filter.MatchedRules(m)
	if len(rules) == 0 {
		log.Debug("%s hit no rules", m.Name)
		return false
	}
	// Check blacklist.
	for _, pattern := range d.cfg.Detector.BlackList {
		ok, err := filepath.Match(pattern, m.Name)
		if err == nil && ok {
			log.Debug("%s hit black pattern %s", m.Name, pattern)
			return false
		}
	}
	return true
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

// store detected metrics.
func (d *Detector) store(m *models.Metric) error {
	// Metric.
	if err := d.db.Metric.Put(m); err != nil {
		return err
	}
	// Index.
	idx := &models.Index{}
	idx.WriteMetric(m)
	if err := d.db.Index.Put(idx); err != nil {
		return err
	}
	return nil
}
