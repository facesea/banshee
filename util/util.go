// Copyright 2015 Eleme Inc. All rights reserved.

package util

import (
	"os"
	"strconv"
	"strings"
)

// ToFixed truncates float64 type to a particular percision in string.
func ToFixed(n float64, prec int) string {
	s := strconv.FormatFloat(n, 'f', prec, 64)
	return strings.TrimRight(s, "0")
}

// Match tests whether a string matches a wildcard pattern.
// Only one or multiple character "*" are supported.
func Match(p, s string) bool {
	l := strings.Split(p, "*")
	if len(l) == 1 {
		return s == p
	}
	for i, j := 0, 0; i < len(l); i++ {
		j = strings.Index(s[j:], l[i])
		if j < 0 {
			return false
		}
		j += len(l[i])
	}
	return true
}

// IsFileExist test whether a filepath is exist.
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
