// Copyright 2015 Eleme Inc. All rights reserved.

package models

import "sync"

// cache is to guarantee goroutine-safety.
type cache struct {
	sync.RWMutex
	// shared is true if the instance is shared betwen goroutines.
	// it shouldn't be reset once set, and shouldn't be added to the global
	// memory until it's set true.
	shared bool
}

// IsShared returns true if the instance is shared between goroutines.
func (c *cache) IsShared() bool {
	return c.shared
}

// MakeShared makes the instance shared as true.
func (c *cache) MakeShared() {
	if !c.IsShared() {
		c.shared = true
	}
}

// Lock if shared.
func (c *cache) LockIfShared() {
	if c.IsShared() {
		c.Lock()
	}
}

// Unlock if shared.
func (c *cache) UnlockIfShared() {
	if c.IsShared() {
		c.Unlock()
	}
}

// Rlock if shared.
func (c *cache) RLockIfShared() {
	if c.IsShared() {
		c.RLock()
	}
}

// RUnlock if shared.
func (c *cache) RUnlockIfShared() {
	if c.IsShared() {
		c.RUnlock()
	}
}
