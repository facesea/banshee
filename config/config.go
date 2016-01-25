// Copyright 2015 Eleme Inc. All rights reserved.

package config

import (
	"encoding/json"
	"github.com/eleme/banshee/util/log"
	"io/ioutil"
)

// Measures
const (
	// Time
	Second uint32 = 1
	Minute        = 60 * Second
	Hour          = 60 * Minute
	Day           = 24 * Hour
)

// Defaults
const (
	// Default time interval for all metrics in seconds.
	DefaultInterval uint32 = 10 * Second
	// Default hit limit to a rule in an interval
	DefaultIntervalHitLimit int = 100
	// Default period for all metrics in seconds.
	DefaultPeriod uint32 = 1 * Day
	// Default weight factor for moving average.
	DefaultWeightFactor float64 = 0.05
	// Default value of alerting interval.
	DefaultAlerterInterval uint32 = 20 * Minute
	// Default value of alert times limit in one day for the same metric
	DefaultAlerterOneDayLimit uint32 = 5
	// Default value of least count.
	DefaultLeastCount uint32 = 5 * Minute / DefaultInterval
)

// Limitations
const (
	// Max value for the number of DefaultTrustLines.
	MaxDefaltTrustLinesLen = 8
	// Max value for the number of FillBlankZeros.
	MaxFillBlankZerosLen = 8
)

// Config is the configuration container.
type Config struct {
	Interval uint32         `json:"interval"`
	Period   uint32         `json:"period"`
	Storage  configStorage  `json:"storage"`
	Detector configDetector `json:"detector"`
	Webapp   configWebapp   `json:"webapp"`
	Alerter  configAlerter  `json:"alerter"`
}

type configStorage struct {
	Path string `json:"path"`
}

type configDetector struct {
	Port              int                `json:"port"`
	Factor            float64            `json:"factor"`
	LeastCount        uint32             `json:"leastCount"`
	BlackList         []string           `json:"blackList"`
	IntervalHitLimit  int                `json:"intervalHitLimit"`
	DefaultTrustLines map[string]float64 `json:"defaultTrustLines"`
	FillBlankZeros    []string           `json:"fillBlankZeros"`
}

type configWebapp struct {
	Port   int       `json:"port"`
	Auth   [2]string `json:"auth"`
	Static string    `json:"static"`
}

type configAlerter struct {
	Command     string `json:"command"`
	Workers     int    `json:"workers"`
	Interval    uint32 `json:"inteval"`
	OneDayLimit uint32 `json:"oneDayLimit"`
}

// New creates a Config with default values.
func New() *Config {
	config := new(Config)
	config.Interval = DefaultInterval
	config.Period = DefaultPeriod
	config.Storage.Path = "storage/"
	config.Detector.Port = 2015
	config.Detector.Factor = DefaultWeightFactor
	config.Detector.LeastCount = DefaultLeastCount
	config.Detector.BlackList = []string{}
	config.Detector.IntervalHitLimit = DefaultIntervalHitLimit
	config.Detector.DefaultTrustLines = make(map[string]float64, 0)
	config.Detector.FillBlankZeros = []string{}
	config.Webapp.Port = 2016
	config.Webapp.Auth = [2]string{"admin", "admin"}
	config.Webapp.Static = "static/dist"
	config.Alerter.Command = ""
	config.Alerter.Workers = 4
	config.Alerter.Interval = DefaultAlerterInterval
	config.Alerter.OneDayLimit = DefaultAlerterOneDayLimit
	return config
}

// UpdateWithJSONFile update the config from a json file.
func (config *Config) UpdateWithJSONFile(fileName string) error {
	log.Debug("read config from %s..", fileName)
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

// Copy config.
func (config *Config) Copy() *Config {
	c := New()
	c.Interval = config.Interval
	c.Period = config.Period
	c.Storage.Path = config.Storage.Path
	c.Detector.Port = config.Detector.Port
	c.Detector.Factor = config.Detector.Factor
	c.Detector.LeastCount = config.Detector.LeastCount
	c.Detector.BlackList = config.Detector.BlackList
	c.Detector.DefaultTrustLines = config.Detector.DefaultTrustLines
	c.Detector.FillBlankZeros = config.Detector.FillBlankZeros
	c.Detector.IntervalHitLimit = config.Detector.IntervalHitLimit
	c.Webapp.Port = config.Webapp.Port
	c.Webapp.Auth = config.Webapp.Auth
	c.Webapp.Static = config.Webapp.Static
	c.Alerter.Command = config.Alerter.Command
	c.Alerter.Workers = config.Alerter.Workers
	c.Alerter.Interval = config.Alerter.Interval
	c.Alerter.OneDayLimit = config.Alerter.OneDayLimit
	return c
}

// Validate config
func (config *Config) Validate() error {
	// Globals
	if config.Interval < 1*Second || config.Interval > 5*Minute {
		return ErrInterval
	}
	if config.Interval > config.Period {
		return ErrPeriod
	}
	// Detector
	if config.Detector.Port < 1 || config.Detector.Port > 65535 {
		return ErrDetectorPort
	}
	if config.Detector.Factor < 0 || config.Detector.Factor > 1 {
		return ErrDetectorFactor
	}
	if len(config.Detector.DefaultTrustLines) > MaxDefaltTrustLinesLen {
		return ErrDetectorDefaultTrustLinesLen
	}
	for _, value := range config.Detector.DefaultTrustLines {
		if value == 0 {
			return ErrDetectorDefaultTrustLineZero
		}
	}
	if len(config.Detector.FillBlankZeros) > MaxFillBlankZerosLen {
		return ErrDetectorFillBlankZerosLen
	}
	// Webapp
	if config.Webapp.Port < 1 || config.Webapp.Port > 65535 {
		return ErrWebappPort
	}
	// Alerter
	if len(config.Alerter.Command) == 0 {
		return ErrAlerterCommandEmpty
	}
	if config.Alerter.Interval <= 0 {
		return ErrAlerterInterval
	}
	if config.Alerter.OneDayLimit <= 0 {
		return ErrAlerterOneDayLimit
	}
	return nil
}
