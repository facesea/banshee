// Copyright 2015 Eleme Inc. All rights reserved.

package metricdb

import "errors"

var (
	// ErrNotFound is returned when requested data not found.
	ErrNotFound = errors.New("metricdb: not found")
	// ErrCorrupted is returned when corrupted data found.
	ErrCorrupted = errors.New("metricdb: corrupted data found")
)
