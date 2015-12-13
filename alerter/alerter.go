// Copyright 2015 Eleme Inc. All rights reserved.

// Alerter will alert via executing command on anomalies detected.
package alerter

import (
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
)

type Alerter struct {
	// Debug
	debug bool
	// Config
	cfg *config.Config
	// Logger
	logger *util.Logger
	// Input
	in chan *models.Metric
}

// New creates a new Alerter.
func New(debug bool, cfg *config.Config, in chan *models.Metric) *Alerter {
	a := new(Alerter)
	a.debug = debug
	a.cfg = cfg
	a.logger = util.NewLogger("banshee.alerter")
	if a.debug {
		a.logger.SetLevel(util.LOG_DEBUG)
	}
	a.in = in
	return a
}

func (a *Alerter) Start() {
	for i := 0; i < a.cfg.Alerter.Workers; i++ {
		// go
	}
}
