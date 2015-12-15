// Copyright 2015 Eleme Inc. All rights reserved.

package adb

import "errors"

var (
	ErrNotFound  = errors.New("adb: not found")
	ErrCorrupted = errors.New("adb: corrupted data found")
)
