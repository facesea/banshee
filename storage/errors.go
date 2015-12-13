// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import "errors"

var (
	ErrNotFound  = errors.New("db: not found")
	ErrCorrupted = errors.New("db: corrupted data found")
)
