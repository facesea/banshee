// Copyright 2015 Eleme Inc. All rights reserved.

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
	"github.com/eleme/banshee/storage/indexdb"
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
	log.Info("detector is listening on tcp://%s..", addr)
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
			log.Info("%dμs detected %s %.4f", elapsed.Nanoseconds()/1000, m.Name, m.Score)
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
	// Bind matched rules.
	m.MatchedRules = rules
	return true
}

// Detect incoming metric with 3-sigma rule and fill the metric.Score.
func (d *Detector) detect(m *models.Metric) error {
	// Get pervious state.
	s, err := d.db.State.Get(m)
	if err != nil && err != statedb.ErrNotFound {
		return err
	}
	if err == statedb.ErrNotFound {
		s = nil
	}
	// Fill blank with zeros.
	for _, pattern := range d.cfg.Detector.FillBlankZeros {
		if ok, _ := filepath.Match(pattern, m.Name); ok {
			// Need to fill zeros.
			if s, err = d.fillBlankZeros(m, s); err != nil {
				return err
			}
			break
		}
	}
	// Move state next.
	s = d.cursor.Next(s, m)
	// Put the next state to db.
	return d.db.State.Put(m, s)
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

// fillBlankZeros moves the state with zero on blanks.
func (d *Detector) fillBlankZeros(m *models.Metric, s *models.State) (*models.State, error) {
	// Get index.
	idx, err := d.db.Index.Get(m.Name)
	if err != nil && err != indexdb.ErrNotFound {
		return nil, err
	}
	if err == indexdb.ErrNotFound {
		// Not found, the metric is the first time seen and s must be nil.
		// The metric will be trusted since its state is nil.
		return s, nil
	}
	// Fill history blanks with zeros.
	numGrid := d.cfg.Period[0]
	gridLen := d.cfg.Period[1]
	period := numGrid * gridLen
	periodNo := m.Stamp / period
	gridNo := (m.Stamp % period) / gridLen
	// Find the start grid stamp
	gridStart := gridNo*gridLen + periodNo*period
	for gridStart-period >= idx.Stamp {
		gridStart -= period
	}
	for ; gridStart < m.Stamp; gridStart += period {
		for stamp := gridStart; stamp < gridStart+gridLen && stamp < m.Stamp; stamp += d.cfg.Interval {
			if stamp > idx.Stamp {
				// Move state with zero.
				s = d.cursor.Next(s, &models.Metric{Stamp: stamp})
			}
		}
	}
	return s, nil
}
