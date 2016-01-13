// Copyright 2015 Eleme Inc. All rights reserved.

// Package util provides util functions.
package util

import (
	"os"
	"strconv"
	"strings"
)

// ToFixed truncates float64 type to a particular percision in string.
func ToFixed(n float64, prec int) string {
	s := strconv.FormatFloat(n, 'f', prec, 64)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

// IsFileExist test whether a filepath is exist.
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
