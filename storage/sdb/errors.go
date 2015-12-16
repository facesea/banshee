// Copyright 2015 Eleme Inc. All rights reserved.

package sdb

import "errors"

var (
	// ErrNotFound is returned when requested data not found.
	ErrNotFound = errors.New("sdb: not found")
	// ErrCorrupted is returned when corrupted data found.
	ErrCorrupted = errors.New("sdb: corrupted data found")
)
