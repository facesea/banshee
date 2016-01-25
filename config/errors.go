package config

import "errors"

// Errors
var (
	// Error
	ErrInterval                     = errors.New("interval should be an integer between 1s~10min")
	ErrPeriod                       = errors.New("period should be an integer greater than interval")
	ErrDetectorPort                 = errors.New("invalid detector.port")
	ErrDetectorFactor               = errors.New("detector.factor should be an float between 0 and 1")
	ErrDetectorDefaultTrustLinesLen = errors.New("detector.defaultTrustLines should have up to 8 items")
	ErrDetectorDefaultTrustLineZero = errors.New("detector.defaultTrustLines should not contain zeros")
	ErrDetectorFillBlankZerosLen    = errors.New("detector.fillBlankZeros should have up to 8 items")
	ErrWebappPort                   = errors.New("invalid webapp.port")
	ErrAlerterInterval              = errors.New("alerter.interval should be greater than 0")
	ErrAlerterOneDayLimit           = errors.New("alerter.oneDayLimit should be greater than 0")
	// Warn
	ErrAlerterCommandEmpty = errors.New("alerter.command is empty")
)
