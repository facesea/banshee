// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import "errors"

var (
	// ErrNotFound is returned when requested project not found.
	ErrNotFound = errors.New("admindb: not found")
	// ErrConstraintUnique is returned when the unique constraint is violated.
	ErrConstraintUnique = errors.New("admindb: unique constraint violated")
)
