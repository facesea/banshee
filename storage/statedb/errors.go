// Copyright 2015 Eleme Inc. All rights reserved.

package statedb

import "errors"

var (
	// ErrNotFound is returned when requested data not found.
	ErrNotFound = errors.New("statedb: not found")
	// ErrCorrupted is returned when corrupted data found.
	ErrCorrupted = errors.New("statedb: corrupted data found")
)
