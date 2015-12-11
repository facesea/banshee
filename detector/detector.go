// Copyright 2015 Eleme Inc. All rights reserved.

//
package detector

import (
	"bufio"
	"fmt"
	"net"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/storage"
)

type Detector struct {
	cfg        *config.Config
	db    *storage.DB
	rules      []string
	rulesCache map[string]bool
	rulesNames map[string][]string
}
var logger = util.NewLogger("detector")
//Init new Detector
func New(cfg *config.Config) *Detector {
	db, err := storage.Open(cfg)
	if err != nil {
		logger.Fatal("failed to open %s: %v", cfg.Storage.Path, err)
	}
	detector := new(Detector)
	detector.cfg = cfg
	detector.db = db
	detector.rulesCache = map[string]bool{};
	detector.rulesNames = map[string][]string{};
	//todo get rules
	return detector
}

//Start detector
func (detector *Detector) Start() {
	addr := fmt.Sprintf("0.0.0.0:%d", detector.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to bind %s: %v", addr, err)
	}
	logger.Info("listening on %s..", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Fatal("failed to accept new conn: %v", err)
		}
		go detector.Handle(conn)
	}
}

func (detector *Detector) Handle(conn net.Conn) {
	addr := conn.RemoteAddr()
	defer func() {
		conn.Close()
		logger.Info("conn %s disconnected", addr)
	}()
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Info("failed to read conn: %v, closing it..", err)
			return
		}
		s := scanner.Text()
		metric, err := parseMetric(s);
		if err != nil {
			logger.Fatal("failed to parse metric: %s", s)
			return
		}
		if detector.Match(metric) {
			detector.Detect(metric)
		}

	}
}

func (detector *Detector) Match(metric *models.Metric) bool {
	bool, exists := detector.rulesCache[metric.Name]
	if bool && exists {
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
			}else {
				detector.rulesNames[pattern] = []string{metric.Name}
			}
			return true
		}
	}
	return false
}

func (detector *Detector) Detect(metric *models.Metric) error {
	//todo
	return nil
}

func (detector *Detector) GetDBData(name string) (avg float64, std float64, num int, err error) {
	//todo
	return
}

func (detector *Detector) PutDBData(key string, avg float64, std float64, num int) error {
	//todo
	return nil
}