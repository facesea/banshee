// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import "errors"

var (
	// ErrProtocol is returned when input line is invalid to parse.
	ErrProtocol = errors.New("detector: invalid protocol input")
	// ErrMetricNameTooLong is returned when input metric name is too long.
	ErrMetricNameTooLong = errors.New("detector: metric name is too long")
	// ErrMetricStampTooSmall is returned when input metric stamp is too small.
	ErrMetricStampTooSmall = errors.New("detector: metric stamp is too small")
)
