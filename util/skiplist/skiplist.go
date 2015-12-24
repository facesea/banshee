// Copyright 2015 Eleme Inc. All rights reserved.

// Package skiplist implements goroutine-safe skiplist with O(1) reads and
// O(logN) writes, mainly for sorted items.
package skiplist

// Compare to SafeMap:
//   Skiplist: Get O(1), Put O(logN), Delete O(log(N)), SortedItems O(N)
//   SafeMap:  Get O(1), Put O(1), Delete O(1), SortedItems O(N*log(N))
// Compare to array with binary-search strageory:
//   the same log(N) search time complexity but less memory moves with
//   put/delete.

import (
	"fmt"
	"math/rand"
	"sync"
)

const LevelMax = 36
const FactorP float64 = 0.5

// Node is skiplist node.
type Node struct {
	score    int
	data     interface{}
	forwards []*Node
}

// Skiplist
type Skiplist struct {
	sync.RWMutex
	length int
	level  int
	head   *Node
	index  map[int]*Node
}

// randLevel returns a level between 1 and LevelMAx.
func randLevel() int {
	level := 1
	for float64(rand.Int()&0xffff) < FactorP*float64(0xffff) {
		level += 1
	}
	if level < LevelMax {
		return level
	}
	return LevelMax
}

// NewNode creates a skiplist node.
func NewNode(level int, score int, data interface{}) *Node {
	return &Node{
		score:    score,
		data:     data,
		forwards: make([]*Node, level),
	}
}

// New creates a skiplist.
func New() *Skiplist {
	return &Skiplist{
		length: 0,
		level:  1,
		head:   NewNode(LevelMax, 0, nil),
		index:  make(map[int]*Node),
	}
}

// Len returns skiplist length.
func (sl *Skiplist) Len() int {
	sl.RLock()
	defer sl.RUnlock()
	return sl.length
}

// Level returns skiplist level.
func (sl *Skiplist) Level() int {
	sl.RLock()
	defer sl.RUnlock()
	return sl.level
}

// Put adds data with score to skiplist. O(logN)
func (sl *Skiplist) Put(score int, data interface{}) {
	sl.Lock()
	defer sl.Unlock()
	update := make([]*Node, LevelMax)
	node := sl.head
	// Find node.
	for i := sl.level - 1; i >= 0; i-- {
		for node.forwards[i] != nil && node.forwards[i].score < score {
			node = node.forwards[i]
		}
		update[i] = node
	}
	// New level.
	level := randLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}
	// Add node.
	node = NewNode(level, score, data)
	for i := 0; i < level; i++ {
		node.forwards[i] = update[i].forwards[i]
		update[i].forwards[i] = node
	}
	// Add to index.
	sl.index[score] = node
	// Incr length.
	sl.length += 1
}

// Get data by score. O(logN)
func (sl *Skiplist) Get(score int) (interface{}, bool) {
	sl.RLock()
	defer sl.RUnlock()
	node, ok := sl.index[score]
	if ok {
		return node.data, true
	}
	return nil, false
}

// Has checks if a score is in list. O(logN)
func (sl *Skiplist) Has(score int) bool {
	sl.RLock()
	defer sl.RUnlock()
	_, ok := sl.index[score]
	return ok
}

// Delete data by score. O(logN)
func (sl *Skiplist) Delete(score int) bool {
	sl.Lock()
	defer sl.Unlock()
	// Find node.
	update := make([]*Node, LevelMax)
	head := sl.head
	node := head
	for i := sl.level - 1; i >= 0; i-- {
		for node.forwards[i] != nil && node.forwards[i].score < score {
			node = node.forwards[i]
		}
		update[i] = node
	}
	node = node.forwards[0]
	if node == nil || node.score != score {
		// Not found.
		return false
	}
	// Delete
	for i := 0; i < sl.level; i++ {
		if update[i].forwards[i] == node {
			update[i].forwards[i] = node.forwards[i]
		}
	}
	// Decr level if need.
	for sl.level > 1 && head.forwards[sl.level-1] == nil {
		sl.level -= 1
	}
	// Delete from index.
	delete(sl.index, score)
	// Decr length.
	sl.length -= 1
	return true
}

// Items returns all sorted items. O(N)
func (sl *Skiplist) Items() []interface{} {
	sl.RLock()
	defer sl.RUnlock()
	items := make([]interface{}, sl.length)
	i := 0
	for node := sl.head.forwards[0]; node != nil; node = node.forwards[0] {
		items[i] = node.data
		i += 1
	}
	return items
}

// Map returns all items in map. O(N)
func (sl *Skiplist) Map() map[int]interface{} {
	sl.RLock()
	defer sl.RUnlock()
	m := make(map[int]interface{}, sl.length)
	for score, node := range sl.index {
		m[score] = node.data
	}
	return m
}

// Print the skiplist.
func (sl *Skiplist) Print() {
	sl.RLock()
	defer sl.RUnlock()
	for i := 0; i < sl.level; i++ {
		node := sl.head.forwards[i]
		fmt.Printf("Level[%d]: ", i)
		for node != nil {
			fmt.Printf("%d[%v] -> ", node.score, node.data)
			node = node.forwards[i]
		}
		fmt.Printf("nil\n")
	}
}
