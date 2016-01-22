// Copyright 2015 Eleme Inc. All rights reserved.

// Package cleaner checks outdated data every certain interval and delete them
// from the storage.
package cleaner

import (
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"time"
)

// The time interval to check outdated data to clean is:
// period / periodNumToInterval. If the period is 1 day, this interval should
// be 4 hours.
const periodNumToInterval uint32 = 6

// The time expiration for metrics is:
// period * metricExpirationNumToPeriod. If the period is 1 day, this
// expiration should be 7 days.
const metricExpirationNumToPeriod uint32 = 7

// The time expiration for metric/index/state is:
// period * expirationNumToPeriod. If the period is 1 day, this expiration
// should be 3 days.
const expirationNumToPeriod uint32 = 3

// Cleaner is to clean outdated data.
type Cleaner struct {
	// Storage
	db *storage.DB
	// Expiration
	metricExpiration uint32
	expiration       uint32
	// Interval
	interval uint32
}

// New creates a cleaner.
func New(db *storage.DB, period uint32) *Cleaner {
	c := new(Cleaner)
	c.db = db
	c.metricExpiration = metricExpirationNumToPeriod * period
	c.expiration = expirationNumToPeriod * period
	c.interval = period / periodNumToInterval
	return c
}

// clean checks all indexes and do cleaning.
func (c *Cleaner) clean() {
	idxs := c.db.Index.All()
	// Use local server time and uint32 is enough for further 90 years
	now := uint32(time.Now().Unix())
	for _, idx := range idxs {
		if idx.Stamp+c.expiration < now {
			// Long time no data, clean all.
			c.db.State.Delete(idx.Name)
			c.db.Index.Delete(idx.Name)
			c.db.Metric.DeleteTo(idx.Name, idx.Stamp+1) // DeleteTo is right closed
			log.Info("%s fully cleaned", idx.Name)
		} else {
			// Clean outdated metrics.
			n, _ := c.db.Metric.DeleteTo(idx.Name, now-c.metricExpiration)
			if n > 0 {
				log.Info("%s %d outdated metrics cleaned", idx.Name, n)
			}
		}
	}
}

// Start a time ticker to clean.
func (c *Cleaner) Start() {
	log.Info("start cleaner with interval %ds..", c.interval)
	// Clean right now.
	c.clean()
	// Clean each interval.
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		<-ticker.C
		c.clean()
	}
}
