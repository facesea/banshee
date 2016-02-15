// Copyright 2015 Eleme Inc. All rights reserved.

// Package cleaner checks outdated data every certain interval and delete them
// from the storage.
package cleaner

import (
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"time"
)

// Cleaner is to clean outdated data.
type Cleaner struct {
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
}

// New creates a cleaner.
func New(cfg *config.Config, db *storage.DB) *Cleaner {
	return &Cleaner{cfg, db}
}

// clean checks all indexes and do cleaning.
func (c *Cleaner) clean() {
	idxs := c.db.Index.All()
	// Use local server time and uint32 is enough for further 90 years
	now := uint32(time.Now().Unix())
	for _, idx := range idxs {
		if idx.Stamp+c.cfg.Cleaner.Threshold < now {
			// Long time no data, clean all.
			c.db.Index.Delete(idx.Name)
			c.db.Metric.DeleteTo(idx.Name, idx.Stamp+1) // DeleteTo is right closed
			log.Info("%s fully cleaned", idx.Name)
		} else {
			// Clean outdated metrics.
			n, _ := c.db.Metric.DeleteTo(idx.Name, now-c.cfg.Expiration)
			if n > 0 {
				log.Debug("%s %d outdated metrics cleaned", idx.Name, n)
			}
		}
	}
}

// Start a time ticker to clean.
func (c *Cleaner) Start() {
	log.Info("start cleaner..")
	// Clean right now.
	c.clean()
	// Clean each interval.
	ticker := time.NewTicker(time.Duration(c.cfg.Cleaner.Interval) * time.Second)
	for {
		<-ticker.C
		c.clean()
	}
}
