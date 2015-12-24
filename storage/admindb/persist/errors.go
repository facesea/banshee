// Copyright 2015 Eleme Inc. All rights reserved.

package persist

import "errors"

var (
	// ErrNotFound is returned when record not found.
	ErrNotFound = errors.New("admin.persist: not found")
	// ErrUnique is returned when the unique constraint is violated.
	ErrUnique = errors.New("admin.persist: unique constraint violated")
	// ErrPrimaryKey is returned when the primary key constraint is violated.
	ErrPrimaryKey = errors.New("admin.persist: primary key constraint violated")
	// ErrNotNull is returned when the not null constraint is violated.
	ErrNotNull = errors.New("admin.persist: not null constraint violated")
)
