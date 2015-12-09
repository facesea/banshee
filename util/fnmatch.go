// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"strings"
)

// FnMatch tests whether a string matches a wildcard pattern.
func FnMatch(s, p string) bool {
	l := strings.Split(p, "*")
	for i, j := 0, 0; i < len(l); i++ {
		j = strings.Index(s[j:], l[i])
		if j < 0 {
			return false
		}
		j += len(l[i])
	}
	return true
}
