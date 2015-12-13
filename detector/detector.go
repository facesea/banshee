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

	"github.com/eleme/banshee/algorithm"
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
	// Rules
	rules      []string
	rulesCache *util.SafeMap
	rulesNames map[string][]string
}

// Init new Detector.
func New(debug bool, cfg *config.Config, db *storage.DB) *Detector {
	detector := new(Detector)
	detector.debug = debug
	detector.cfg = cfg
	detector.logger = util.NewLogger("detector")
	if detector.debug {
		detector.logger.SetLevel(util.LOG_DEBUG)
	}
	detector.db = db
	detector.rulesCache = util.NewSafeMap()
	detector.rulesNames = map[string][]string{}
	// FIXME: rules
	return detector
}

// Start detector
func (detector *Detector) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", detector.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		detector.logger.Fatal("failed to bind tcp://%s: %v", addr, err)
	}
	detector.logger.Info("listening on tcp://%s..", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			detector.logger.Fatal("failed to accept new conn: %v", err)
		}
		go detector.handle(conn)
	}
}

// Handle a connection, it will filter the mertics by rules and detect whether
// the metrics are anomalies.
func (detector *Detector) handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	defer func() {
		conn.Close()
		detector.logger.Info("conn %s disconnected", addr)
	}()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			detector.logger.Info("failed to read conn: %v, closing it..", err)
			break
		}
		line := scanner.Text()
		m, err := parseMetric(line)
		if err != nil {
			if len(line) > 10 {
				line = line[:10]
			}
			detector.logger.Error("failed to parse '%s': %v, skipping..", line, err)
			continue
		}
		if detector.match(m) {
			err = detector.detect(m)
			if err != nil {
				detector.logger.Error("failed to detect metric: %v, skipping..", err)
				continue
			}
			detector.logger.Debug("detected %s => average %.3f, socre %.3f", m.Name, m.Average, m.Score)
		}
	}
}

func (detector *Detector) match(m *models.Metric) bool {
	// FIXME
	// v, ok := detector.rulesCache.Get(m.Name)
	// b := v.(bool)
	// if b && ok {
	// 	return true
	// }

	for _, pattern := range detector.cfg.Detector.BlackList {
		matched := util.FnMatch(pattern, m.Name)
		if matched {
			return false
		}
	}
	return true // FIXME: return true tempory
	// FIXME: get rules from db
	for _, pattern := range detector.rules {
		matched := util.FnMatch(pattern, m.Name)
		if matched {
			detector.rulesCache.Set(m.Name, true)
			slice, exists := detector.rulesNames[pattern]
			if exists {
				detector.rulesNames[pattern] = append(slice, m.Name)
			} else {
				detector.rulesNames[pattern] = []string{m.Name}
			}
			return true
		}
	}
	return false
}

// Detect incoming metric with 3-sigma rule and fill the metric.Score.
func (detector *Detector) detect(m *models.Metric) error {
	wf := detector.cfg.Detector.TrendFactor
	startSize := detector.cfg.Detector.StartSize
	state, err := detector.db.GetState(m)
	if err != nil && err != storage.ErrNotFound {
		return err
	}
	stateNext := &models.State{}
	if err == storage.ErrNotFound {
		m.Average = m.Value
		stateNext.Average = m.Value
		stateNext.StdDev = 0
		stateNext.Count = 1
	} else {
		m.Average = state.Average
		stateNext.Average = algorithm.Ewma(wf, state.Average, m.Value)
		stateNext.StdDev = algorithm.Ewms(wf, state.Average, stateNext.Average, state.StdDev, m.Value)
		if state.Count < startSize {
			stateNext.Count = state.Count + 1
		} else {
			stateNext.Count = state.Count
		}
	}
	if stateNext.Count >= startSize {
		m.Score = algorithm.Div3Sigma(stateNext.Average, stateNext.StdDev, m.Value)
	} else {
		m.Score = 0
	}
	err = detector.db.PutState(m, stateNext)
	return err
}
