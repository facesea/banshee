// Copyright 2015 Eleme Inc. All rights reserved.

package models

import (
	"errors"
	"sync"
)

type cache struct {
	lock sync.RWMutex
	// shared is readonly once set
	shared bool
}

// Share the instance between goroutines.
func (c *cache) Share() {
	if c.shared {
		panic(errors.New("cache: instance already shared"))
	}
	c.shared = true
}

// Lock locks instance for writing.
func (c *cache) Lock() {
	if c.shared {
		c.lock.Lock()
	}
}

// Unlock unlocks instance writing.
func (c *cache) Unlock() {
	if c.shared {
		c.lock.Unlock()
	}
}

// RLock locks instance for reading.
func (c *cache) RLock() {
	if c.shared {
		c.lock.RLock()
	}
}

// RUnlock unlocks instance reading.
func (c *cache) RUnlock() {
	if c.shared {
		c.lock.RUnlock()
	}
}
