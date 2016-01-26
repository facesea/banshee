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
	// Default metric expiration.
	DefaultExpiration uint32 = 7 * Day
	// Default weight factor for moving average.
	DefaultTrendingFactor float64 = 0.05
	// Default filter offset to query history metrics.
	DefaultFilterOffset float64 = 0.01
	// Default cleaner interval.
	DefaultCleanerInterval uint32 = 3 * Hour
	// Default cleaner threshold.
	DefaultCleanerThreshold uint32 = 3 * Day
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
	// Min value for the expiration to period.
	MinExpirationNumToPeriod uint32 = 5
	// Min value for the cleaner threshold to period.
	MinCleanerThresholdNumToPeriod uint32 = 2
)

// Config is the configuration container.
type Config struct {
	Interval   uint32         `json:"interval"`
	Period     uint32         `json:"period"`
	Expiration uint32         `json:"expiration"`
	Storage    configStorage  `json:"storage"`
	Detector   configDetector `json:"detector"`
	Webapp     configWebapp   `json:"webapp"`
	Alerter    configAlerter  `json:"alerter"`
	Cleaner    configCleaner  `json:"cleaner"`
}

type configStorage struct {
	Path string `json:"path"`
}

type configDetector struct {
	Port              int                `json:"port"`
	TrendingFactor    float64            `json:"trendingFactor"`
	FilterOffset      float64            `json:"filterOffset"`
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

type configCleaner struct {
	Interval  uint32 `json:"interval"`
	Threshold uint32 `json:"threshold"`
}

// New creates a Config with default values.
func New() *Config {
	config := new(Config)
	config.Interval = DefaultInterval
	config.Period = DefaultPeriod
	config.Expiration = DefaultExpiration
	config.Storage.Path = "storage/"
	config.Detector.Port = 2015
	config.Detector.TrendingFactor = DefaultTrendingFactor
	config.Detector.FilterOffset = DefaultFilterOffset
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
	config.Cleaner.Interval = DefaultCleanerInterval
	config.Cleaner.Threshold = DefaultCleanerThreshold
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
	c.Expiration = config.Expiration
	c.Storage.Path = config.Storage.Path
	c.Detector.Port = config.Detector.Port
	c.Detector.TrendingFactor = config.Detector.TrendingFactor
	c.Detector.FilterOffset = config.Detector.FilterOffset
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
	c.Cleaner.Interval = config.Cleaner.Interval
	c.Cleaner.Threshold = config.Cleaner.Threshold
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
	if config.Expiration < config.Period*MinExpirationNumToPeriod {
		return ErrExpiration
	}
	// Detector
	if config.Detector.Port < 1 || config.Detector.Port > 65535 {
		return ErrDetectorPort
	}
	if config.Detector.TrendingFactor <= 0 || config.Detector.TrendingFactor >= 1 {
		return ErrDetectorTrendingFactor
	}
	if config.Detector.FilterOffset <= 0 || config.Detector.FilterOffset >= 1 {
		return ErrDetectorFilterOffset
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
	// Cleaner
	if config.Cleaner.Threshold < config.Period*MinCleanerThresholdNumToPeriod {
		return ErrCleanerThreshold
	}
	return nil
}
