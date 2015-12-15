// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"os"
	"strings"
)

// Strings

// Match tests whether a string matches a wildcard pattern.
// Only one or multiple character "*" are supported.
func Match(s, p string) bool {
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

// File system

// IsFileExist test whether a filepath is exist.
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
