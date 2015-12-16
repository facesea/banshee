// Copyright 2015 Eleme Inc. All rights reserved.

package adb

import "errors"

var (
	// ErrNotFound is returned when requested data not found.
	ErrNotFound = errors.New("adb: not found")
	// ErrCorrupted is returned when corrupted data found.
	ErrCorrupted = errors.New("adb: corrupted data found")
)
