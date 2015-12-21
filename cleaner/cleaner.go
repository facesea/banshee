// Copyright 2015 Eleme Inc. All rights reserved.

// Package cleaner checks outdated data every certain interval and delete them
// from the storage.
package cleaner

import (
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"time"
)

// The time interval to check data to clean is this value * period.
// If the period is 1 day, this interval will be about 4 hours.
const intervalNumToPeriod float32 = 0.2

// The expiration time is this value * period.
// If the period is 1 day, this expiration will be 7 days.
const expirationNumToPeriod float32 = 7

// Cleaner is to help clean outdated data.
type Cleaner struct {
	// DB
	db *storage.DB
	// Expiration
	expiration time.Duration
	// Ticker
	ticker *time.Ticker
	// Interval
	interval time.Duration
}

// New creates a cleaner.
func New(db *storage.DB, period int) *Cleaner {
	c := new(Cleaner)
	c.db = db
	c.expiration = time.Duration(uint32(expirationNumToPeriod*float32(period))) * time.Second
	c.interval = time.Duration(uint32(intervalNumToPeriod*float32(period))) * time.Second
	c.ticker = time.NewTicker(c.interval)
	return c
}

// Start a time ticker and wait to check.
func (c *Cleaner) Start() {
	// Check right now.
	c.clean()
	for {
		// And wait for another interval to check.
		<-c.ticker.C
		c.clean()
	}
}

// clean checks all indexes for outdated metrics, states and clean them.
func (c *Cleaner) clean() {
	idxs := c.db.Index.All()
	now := time.Now()
	for _, idx := range idxs {
		t := time.Unix(int64(idx.Stamp), 0)
		if t.Add(c.expiration).Before(now) {
			// Clean outdated.
			c.db.State.Delete(idx.Name)
			c.db.Metric.DeleteTo(idx.Name, uint32(now.Unix()))
			c.db.Index.Delete(idx.Name)
			log.Info("%s cleaned", idx.Name)
		}
	}
}
