// Copyright 2015 Eleme Inc. All rights reserved.

package admindb

import "errors"

var (
	// ErrNotFound is returned when requested data not found.
	ErrNotFound = errors.New("admindb: not found")
	// ErrProjectNotFound is returned when requested project not found.
	ErrProjectNotFound = errors.New("admindb: project not found")
	// ErrUserNotFound is returned when requested user not found.
	ErrUserNotFound = errors.New("admindb: user not found")
	// ErrRuleNotFound is returned when requested rule not found.
	ErrRuleNotFound = errors.New("admindb: rule not found")
	// ErrConstraintUnique is returned when the unique constraint is violated.
	ErrConstraintUnique = errors.New("admindb: unique constraint violated")
	// ErrConstraintPrimaryKey is returned when the primary key constraint is
	// violated.
	ErrConstraintPrimaryKey = errors.New("admindb: primary key constraint violated")
	// ErrConstraintNotNull is returned when the not null constraint is
	// violated.
	ErrConstraintNotNull = errors.New("admindb: not null constraint violated")
)
