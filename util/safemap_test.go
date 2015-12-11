// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"fmt"
	"testing"
)

func TestSafeMap(t *testing.T) {
	m := NewSafeMap()
	m.Set("key", "val")
	Assert(t, m.Len() == 1)
	val, ok := m.Get("key")
	Assert(t, ok)
	Assert(t, val == "val")
	for k, v := range m.Items() {
		fmt.Printf("%v=>%v", k, v)
	}
	Assert(t, m.Items()["key"] == "val")
	Assert(t, m.Delete("key"))
	Assert(t, !m.Delete("key1"))
}
