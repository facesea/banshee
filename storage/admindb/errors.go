// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import "errors"

var (
	// ErrNotFound is returned when requested data not found.
	ErrNotFound = errors.New("admindb: not found")
	// ErrCorrupted is returned when corrupted data found.
	ErrCorrupted = errors.New("admindb: corrupted data found")
)
