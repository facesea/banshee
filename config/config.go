// Copyright 2015 Eleme Inc. All rights reserved.

// Package config handles the configuration container and parsing.
package config

import (
	"encoding/json"
	"io/ioutil"
)

// Defaults
const (
	// Default time interval for all metrics in seconds.
	DefaultInterval int = 10
	// Default number of grids in one period.
	DefaultNumGrid int = 288
	// Default grid length in seconds.
	DefaultGridLen int = 300
	// Default weight factor for moving average and standard deviation.
	DefaultWeightFactor float64 = 0.05
)

// Least count.
const (
	// Percentage the leastC in one grid.
	leastCountGridPercent float64 = 0.6
	// Min value of leastC.
	leastCountMin int = 18
)

// Config is the configuration container.
type Config struct {
	Interval int            `json:"interval"`
	Period   [2]int         `json:"period"`
	Storage  configStorage  `json:"storage"`
	Detector configDetector `json:"detector"`
	Webapp   configWebapp   `json:"webapp"`
	Alerter  configAlerter  `json:"alerter"`
}

type configStorage struct {
	Path string `json:"path"`
}

type configDetector struct {
	Port      int      `json:"port"`
	Factor    float64  `json:"factor"`
	BlackList []string `json:"blackList"`
}

type configWebapp struct {
	Port int       `json:"port"`
	Auth [2]string `json:"auth"`
}

type configAlerter struct {
	Command string `json:"command"`
	Workers int    `json:"workers"`
}

// New creates a Config with default values.
func New() *Config {
	config := new(Config)
	config.Interval = DefaultInterval
	config.Period = [2]int{DefaultNumGrid, DefaultGridLen}
	config.Storage.Path = "storage/"
	config.Detector.Port = 2015
	config.Detector.Factor = DefaultWeightFactor
	config.Detector.BlackList = []string{}
	config.Webapp.Port = 2016
	config.Webapp.Auth = [2]string{"admin", "admin"}
	config.Alerter.Command = ""
	config.Alerter.Workers = 4
	return config
}

// UpdateWithJSONFile update the config from a json file.
func (config *Config) UpdateWithJSONFile(fileName string) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}
	return err
}

// LeastC returns the least count to start detection, if the count of a metric
// is less than this value, it will be trusted without an calculation on its
// score.
func (config *Config) LeastC() int {
	c := int((float64(config.Period[1]) / float64(config.Interval)) * leastCountGridPercent)
	if c > leastCountMin {
		return c
	}
	return leastCountMin
}
