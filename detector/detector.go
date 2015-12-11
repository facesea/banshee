// Copyright 2015 Eleme Inc. All rights reserved.

// Detector is a tcp server to detect anomalies.
//   detector := New(cfg, db)
//   detector.Start()
//
package detector

import (
	"bufio"
	"fmt"
	"net"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util"
)

// Detector is a tcp server to detect anomalies.
type Detector struct {
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
	// Rules
	rules      []string
	rulesCache map[string]bool
	rulesNames map[string][]string
}

var logger = util.NewLogger("detector")

// Init new Detector.
func New(cfg *config.Config, db *storage.DB) *Detector {
	detector := new(Detector)
	detector.cfg = cfg
	detector.db = db
	detector.rulesCache = map[string]bool{}
	detector.rulesNames = map[string][]string{}
	// FIXME: rules
	return detector
}

// Start detector
func (detector *Detector) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", detector.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to bind tcp://%s: %v", addr, err)
	}
	logger.Info("listening on tcp://%s..", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Fatal("failed to accept new conn: %v", err)
		}
		go detector.handle(conn)
	}
}

func (detector *Detector) handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	defer func() {
		conn.Close()
		logger.Info("conn %s disconnected", addr)
	}()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Info("failed to read conn: %v, closing it..", err)
			break
		}
		line := scanner.Text()
		metric, err := parseMetric(line)
		if err != nil {
			if len(line) > 10 {
				line = line[:10]
			}
			logger.Error("failed to parse '%s': %v, skipping..", line, err)
			continue
		}
		if detector.match(metric) {
			detector.detect(metric)
		}
	}
}

func (detector *Detector) match(metric *models.Metric) bool {
	b, ok := detector.rulesCache[metric.Name]
	if b && ok {
		return true
	}

	for _, pattern := range detector.cfg.Detector.BlackList {
		matched := util.FnMatch(pattern, metric.Name)
		if matched {
			return false
		}
	}

	for _, pattern := range detector.rules {
		matched := util.FnMatch(pattern, metric.Name)
		if matched {
			detector.rulesCache[metric.Name] = true
			slice, exists := detector.rulesNames[pattern]
			if exists {
				detector.rulesNames[pattern] = append(slice, metric.Name)
			} else {
				detector.rulesNames[pattern] = []string{metric.Name}
			}
			return true
		}
	}
	return false
}

func (detector *Detector) detect(metric *models.Metric) error {
	// TODO
	return nil
}
