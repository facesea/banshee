// Copyright 2015 Eleme Inc. All rights reserved.

// Configuration for banshee with default values.
//  Global Options
//    debug            if on debug mode. [default: false]
//    interval         incomding metrics time interval (in sec). [default: 10]
//    periodicity      metrics periodicity (in sec), NumTimeSpans x TimeSpan.
//                     [default: [480, 180]]
//  Storage Options
//    path             path for leveldb, which maintains analyzation
//                     results. [default: "storage/"]
//  Detector Options
//    port             detector tcp port to listen. [default: 2015]
//    trendFactor      the factor to calculate trending value via weighted
//                     moving average algorithm. [default: 0.07]
//    strict           if this is set false, detector will passivate latest
//                     metric. [default: true]
//    blacklist        metrics blacklist, detector will allow one metric to pass
//                     only if it matches one pattern in rules and dosent match
//                     any pattern in blacklist. [default: ["statsd.*"]]
//    startSize        detector won't start to detect until the data set is
//                     larger than this size. [default: 32]
//  WebApp Options
//    port             webapp http port to listen. [default: 2016]
//    auth             username and password for admin basic auth. [default:
//                     ["admin", "admin"]]
//  Alerter Options
//    command          shell command to execute on anomalies detected, leaving
//                     empty means do nothing. [default: ""]
// See also exampleConfig.json please.
//
package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Debug       bool           `json:"debug"`
	Interval    int            `json:"interval"`
	Periodicity [2]int         `json:"periodicity"`
	Storage     ConfigStorage  `json:"storage"`
	Detector    ConfigDetector `json:"detector"`
	Webapp      ConfigWebapp   `json:"webapp"`
	Alerter     ConfigAlerter  `json:"alerter"`
}

type ConfigStorage struct {
	Path string `json:"path"`
}

type ConfigDetector struct {
	Port        int      `json:"port"`
	TrendFactor float64  `json:"trendFactor"`
	Strict      bool     `json:"strict"`
	BlackList   []string `json:"blackList"`
	StartSize   int      `json:"startSize"`
}

type ConfigWebapp struct {
	Port int       `json:"port"`
	Auth [2]string `json:"auth"`
}

type ConfigAlerter struct {
	Command string `json:"command"`
}

// NewConfigWithDefaults creates a Config with default values.
func NewConfigWithDefaults() *Config {
	config := new(Config)
	config.Debug = false
	config.Interval = 10
	config.Periodicity = [2]int{480, 180}
	config.Storage.Path = "storage/"
	config.Detector.Port = 2015
	config.Detector.TrendFactor = 0.07
	config.Detector.Strict = true
	config.Detector.BlackList = []string{"statsd.*"}
	config.Detector.StartSize = 32
	config.Webapp.Port = 2016
	config.Webapp.Auth = [2]string{"admin", "admin"}
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
