// Copyright 2015 Eleme Inc. All rights reserved.

package sdb

import "errors"

var (
	ErrNotFound  = errors.New("sdb: not found")
	ErrCorrupted = errors.New("sdb: corrupted data found")
)
