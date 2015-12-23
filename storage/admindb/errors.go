// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import "errors"

var (
	// ErrNotFound is returned when requested project not found.
	ErrNotFound = errors.New("admindb: not found")
	// ErrConstraintUnique is returned when the unique constraint is violated.
	ErrConstraintUnique = errors.New("admindb: unique constraint violated")
	// ErrConstraintPrimaryKey is returned when the primary key constraint is
	// violated.
	ErrConstraintPrimaryKey = errors.New("admindb: primary key constraint violated")
	// ErrConstraintNotNull is returned when the not null constraint is
	// violated.
	ErrConstraintNotNull = errors.New("admindb: not null constraint violated")
)
