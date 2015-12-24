// Copyright 2015 Eleme Inc. All rights reserved.

package cache

import "errors"

var (
	// ErrProjectNotFound is returned when requested project not found.
	ErrProjectNotFound = errors.New("admindb.cache: project not found")
	// ErrUserNotFound is returned when requested user not found.
	ErrUserNotFound = errors.New("admindb.cache: user not found")
	// ErrRuleNotFound is returned when requested rule not found.
	ErrRuleNotFound = errors.New("admindb.cache: rule not found")
)
