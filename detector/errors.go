// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import "errors"

var (
	// ErrProtocol is returned when input line is invalid to parse.
	ErrProtocol = errors.New("detector: invalid protocol input")
)
