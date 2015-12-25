// Copyright 2015 Eleme Inc. All rights reserved.

// Package safemap implements a map like container with rw lock to keep
// groutine safety.
package safemap

import "sync"

// SafeMap is a type like map with rw lock.
type SafeMap struct {
	lock *sync.RWMutex
	m    map[interface{}]interface{}
}

// New creates a SafeMap.
func New() *SafeMap {
	return &SafeMap{
		lock: &sync.RWMutex{},
		m:    make(map[interface{}]interface{}),
	}
}

// Get value from map by key.
func (m *SafeMap) Get(key interface{}) (interface{}, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.m[key]; ok {
		return val, true
	}
	return nil, false
}

// Set value to map by key.
func (m *SafeMap) Set(key, val interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.m[key] = val
}

// Has checks if a key is in map.
func (m *SafeMap) Has(key interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.m[key]
	return ok
}

// Delete a value from map by key.
func (m *SafeMap) Delete(key interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return true
	}
	return false
}

// Pop a value from map by key.
func (m *SafeMap) Pop(key interface{}) (interface{}, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return v, true
	}
	return nil, false
}

// Items returns all items in the map.
func (m *SafeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	d := make(map[interface{}]interface{}, len(m.m))
	for key, val := range m.m {
		d[key] = val
	}
	return d
}

// Len returns map length.
func (m *SafeMap) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.m)
}

// Clear the map.
func (m *SafeMap) Clear() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Rely on GC
	m.m = make(map[interface{}]interface{})
}
