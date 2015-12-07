// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Interval    int            `json:"interval"`
	Periodicity int            `json:"periodicity"`
	Storage     ConfigStorage  `json:"storage"`
	Detector    ConfigDetector `json:"detector"`
	Webapp      ConfigWebapp   `json:"webapp"`
	Alerter     ConfigAlerter  `json:"alerter"`
}

type ConfigStorage struct {
	Stats string `json:"stats"`
	Rules string `json:"rules"`
}

type ConfigDetector struct {
	Debug       bool     `json:"debug"`
	Port        int      `json:"port"`
	TrendFactor float64  `json:"trendFactor"`
	Strict      bool     `json:"strict"`
	WhiteList   []string `json:"whitelist"`
	BlackList   []string `json:"blackList"`
	StartSize   int      `json:"startSize"`
}

type ConfigWebapp struct {
	Port int    `json:"port"`
	Auth string `json:"auth"`
}

type ConfigAlerter struct {
	Command string `json:"command"`
}

// NewConfigWithDefaults creates a Config with default values.
func NewConfigWithDefaults() *Config {
	config := new(Config)
	config.Interval = 10
	config.Periodicity = 24 * 60 * 60
	config.Storage.Stats = "stats.db"
	config.Storage.Rules = "rules.db"
	config.Detector.Debug = false
	config.Detector.Port = 2015
	config.Detector.TrendFactor = 0.07
	config.Detector.Strict = true
	config.Detector.WhiteList = []string{"*"}
	config.Detector.BlackList = []string{"statsd.*"}
	config.Detector.StartSize = 32
	config.Webapp.Port = 2016
	config.Webapp.Auth = "admin:admin"
	config.Alerter.Command = ""
	return config
}

// NewConfigWithJsonBytes creates a Config with json literal bytes.
func NewConfigWithJsonBytes(b []byte) (*Config, error) {
	config := NewConfigWithDefaults()
	err := json.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// NewConfigWithJsonFile creates a Config from a json file by fileName.
func NewConfigWithJsonFile(fileName string) (*Config, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return NewConfigWithJsonBytes(b)
}
