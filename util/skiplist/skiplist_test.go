// Copyright 2015 Eleme Inc. All rights reserved.

package skiplist

import (
	"github.com/eleme/banshee/util/assert"
	"github.com/eleme/banshee/util/safemap"
	"math/rand"
	"sort"
	"testing"
)

// Bench
var sl *Skiplist
var N int
var m *safemap.SafeMap

func init() {
	N = 1024
	sl = New()
	m = safemap.New()
	for i := 0; i < N; i++ {
		sl.Put(i, rand.Int()&N)
		m.Set(i, rand.Int()&N)
	}
}

// TestBasic
func TestBasic(t *testing.T) {
	sl := New()
	// Suite
	a, b, c, d, e, f, g, h, i := 4, 2, 3, 7, 5, 6, 9, 8, 1
	// Put
	sl.Put(a, "a")
	sl.Put(b, "b")
	sl.Put(c, "c")
	sl.Put(d, "d")
	sl.Put(e, "e")
	sl.Put(f, "f")
	sl.Put(g, "g")
	sl.Put(h, "h")
	sl.Put(i, "i")
	// Get
	v, ok := sl.Get(a)
	assert.Ok(t, ok && "a" == v.(string))
	v, ok = sl.Get(b)
	assert.Ok(t, ok && "b" == v.(string))
	v, ok = sl.Get(c)
	assert.Ok(t, ok && "c" == v.(string))
	v, ok = sl.Get(d)
	assert.Ok(t, ok && "d" == v.(string))
	v, ok = sl.Get(e)
	assert.Ok(t, ok && "e" == v.(string))
	v, ok = sl.Get(f)
	assert.Ok(t, ok && "f" == v.(string))
	v, ok = sl.Get(g)
	assert.Ok(t, ok && "g" == v.(string))
	v, ok = sl.Get(h)
	assert.Ok(t, ok && "h" == v.(string))
	v, ok = sl.Get(i)
	assert.Ok(t, ok && "i" == v.(string))
	// Len
	assert.Ok(t, 9 == sl.Len())
	// Items
	items := sl.Items()
	assert.Ok(t, len(items) == 9)
	assert.Ok(t, "i" == items[0].(string))
	assert.Ok(t, "b" == items[1].(string))
	assert.Ok(t, "c" == items[2].(string))
	assert.Ok(t, "a" == items[3].(string))
	assert.Ok(t, "e" == items[4].(string))
	assert.Ok(t, "f" == items[5].(string))
	assert.Ok(t, "d" == items[6].(string))
	assert.Ok(t, "h" == items[7].(string))
	assert.Ok(t, "g" == items[8].(string))
	// Map
	m := sl.Map()
	x := map[int]string{a: "a", b: "b", c: "c", d: "d", e: "e", f: "f", g: "g", h: "h", i: "i"}
	assert.Ok(t, len(m) == len(x))
	for k, v := range m {
		assert.Ok(t, m[k] == v.(string))
	}
	// Delete
	assert.Ok(t, sl.Delete(d))
	assert.Ok(t, 8 == sl.Len())
	assert.Ok(t, !sl.Has(d))
	assert.Ok(t, sl.Delete(g))
	assert.Ok(t, 7 == sl.Len())
	assert.Ok(t, !sl.Has(g))
	// Print
	sl.Print()
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sl.Get(i & N)
	}
}

func BenchmarkSafeMapGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Get(i & N)
	}
}

func BenchmarkPut(b *testing.B) {
	sl := New()
	for i := 0; i < b.N; i++ {
		sl.Put(i, i)
	}
}

func BenchmarkSafeMapPut(b *testing.B) {
	m := safemap.New()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

func BenchmarkSortedItems(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Sorted items
		sl.Items()
	}
}
func BenchmarkSafeMapSortedItems(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Sort map keys
		d := m.Items()
		ints := make([]int, len(d))
		j := 0
		for k, _ := range d {
			ints[j] = k.(int)
			j++
		}
		sort.Ints(ints)
	}
}
