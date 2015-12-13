// Copyright 2015 Eleme Inc. All rights reserved.

// Detector is a tcp server to detect anomalies.
//
//   detector := New(cfg, db)
//   detector.Start()
//
package detector

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util"
)

// Detector is a tcp server to detect anomalies.
type Detector struct {
	// Debug
	debug bool
	// Config
	cfg *config.Config
	// Logger
	logger *util.Logger
	// Storage
	db *storage.DB
	// Output
	out chan *models.Metric
	// Rules
	rules      []string
	rulesCache *util.SafeMap
	rulesNames map[string][]string
}

// Init new Detector.
func New(debug bool, cfg *config.Config, db *storage.DB, out chan *models.Metric) *Detector {
	d := new(Detector)
	d.debug = debug
	d.cfg = cfg
	d.logger = util.NewLogger("banshee.detector")
	if d.debug {
		d.logger.SetLevel(util.LOG_DEBUG)
	}
	d.db = db
	d.out = out
	d.rulesCache = util.NewSafeMap()
	d.rulesNames = map[string][]string{}
	// FIXME: rules
	return d
}

// Start detector
func (d *Detector) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", d.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		d.logger.Fatal("failed to bind tcp://%s: %v", addr, err)
	}
	d.logger.Info("listening on tcp://%s..", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			d.logger.Fatal("failed to accept new conn: %v", err)
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
		d.logger.Info("conn %s disconnected", addr)
	}()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			d.logger.Info("failed to read conn: %v, closing it..", err)
			break
		}
		startAt := time.Now()
		line := scanner.Text()
		m, err := parseMetric(line)
		if err != nil {
			if len(line) > 10 {
				line = line[:10]
			}
			d.logger.Error("failed to parse '%s': %v, skipping..", line, err)
			continue
		}
		if d.match(m) {
			err = d.detect(m)
			if err != nil {
				d.logger.Error("failed to detect metric: %v, skipping..", err)
				continue
			}
			elapsed := time.Since(startAt)
			d.logger.Debug("detected %s cost=%dμs", m.String(), elapsed.Nanoseconds()/1000)
			d.out <- m
		}
	}
}

func (d *Detector) match(m *models.Metric) bool {
	// FIXME
	// v, ok := d.rulesCache.Get(m.Name)
	// b := v.(bool)
	// if b && ok {
	// 	return true
	// }

	for _, pattern := range d.cfg.Detector.BlackList {
		matched := util.FnMatch(pattern, m.Name)
		if matched {
			return false
		}
	}
	return true // FIXME: return true tempory
	// FIXME: get rules from db
	for _, pattern := range d.rules {
		matched := util.FnMatch(pattern, m.Name)
		if matched {
			d.rulesCache.Set(m.Name, true)
			slice, exists := d.rulesNames[pattern]
			if exists {
				d.rulesNames[pattern] = append(slice, m.Name)
			} else {
				d.rulesNames[pattern] = []string{m.Name}
			}
			return true
		}
	}
	return false
}

// Detect incoming metric with 3-sigma rule and fill the metric.Score.
func (d *Detector) detect(m *models.Metric) error {
	wf := d.cfg.Detector.TrendFactor
	startSize := d.cfg.Detector.StartSize
	state, err := d.db.GetState(m)
	// Unexcepted error
	if err != nil && err != storage.ErrNotFound {
		return err
	}
	stateN := &models.State{}
	if err == storage.ErrNotFound {
		// Not found, initialize as first
		m.Average = m.Value
		stateN.Average = m.Value
		stateN.StdDev = 0
		stateN.Count = 1
	} else {
		// Found, move to next
		m.Average = state.Average
		stateN.Average = ewma(wf, state.Average, m.Value)
		stateN.StdDev = ewms(wf, state.Average, stateN.Average, state.StdDev, m.Value)
		if state.Count < startSize {
			stateN.Count = state.Count + 1
		} else {
			stateN.Count = state.Count
		}
	}
	// Don't calculate the score if current count is not enough.
	if stateN.Count >= startSize {
		m.Score = div3Sigma(stateN.Average, stateN.StdDev, m.Value)
	} else {
		m.Score = 0
	}
	err = d.db.PutState(m, stateN)
	return err
}
