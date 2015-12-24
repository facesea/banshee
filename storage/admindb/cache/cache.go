// Copyright 2015 Eleme Inc. All rights reserved.

// Package cache handles admin cache.
package cache

import (
	"github.com/eleme/banshee/util/skiplist"
)

// Cache.
type Cache struct {
	projs *skiplist.Skiplist
	rules *skiplist.Skiplist
	users *skiplist.Skiplist
}

// New creates cache.
func New() *Cache {
	c := new(Cache)
	c.projs = skiplist.New()
	c.rules = skiplist.New()
	c.users = skiplist.New()
	return c
}
