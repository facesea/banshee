// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"github.com/eleme/banshee/util/log"
	"time"
)

// The time interval to check data to clean is this value * period.
const cleanIntervalNumToPeriod float32 = 0.2

// The expiration time is this value * period.
const expirationNumToPeriod = 7

// cleaner checks outdated data every certain interval and delete them from
// the storage, the targets to be checked are outdated indexes, metrics and
// states.
type cleaner struct {
	// DB
	db *DB
	// Period
	period int
}

// start the time ticker waiting to work.
func (c *cleaner) start() {
	interval := time.Duration(cleanIntervalNumToPeriod * float32(c.period))
	ticker := time.NewTicker(interval * time.Second)
	// Work right now.
	c.work()
	for {
		// And wait for another interval to work.
		<-ticker.C
		c.work()
	}
}

// work checks the index for outdated metrics and states and clean them
// including the index.
func (c *cleaner) work() {
	idxs := c.db.Index.All()
	exp := uint32(expirationNumToPeriod * c.period)
	now := uint32(time.Now().Unix())
	for _, idx := range idxs {
		if idx.Stamp+exp < now {
			c.db.State.Delete(idx.Name)
			c.db.Metric.DeleteTo(idx.Name, now)
			c.db.Index.Delete(idx.Name)
			log.Info("clean %s from db..", idx.Name)
		}
	}
}
