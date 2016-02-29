// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"errors"
	"regexp"
	"strings"
)

// Limitations
const (
	// Max value of the project name length.
	MaxProjectNameLen = 64
	// Max value of the user name length.
	MaxUserNameLen = 32
	// Max value of the rule pattern length.
	MaxRulePatternLen = 256
	// Max value of the metric name length.
	MaxMetricNameLen = 256
	// Min value of the metric stamp.
	MinMetricStamp uint32 = 1450322633
)

// Errors
var (
	ErrProjectNameEmpty         = errors.New("project name is empty")
	ErrProjectNameTooLong       = errors.New("project name is too long")
	ErrProjectSilentTimeStart   = errors.New("project silent time start is invalid")
	ErrProjectSilentTimeEnd     = errors.New("project silent time end is invalid")
	ErrProjectSilentTimeRange   = errors.New("project silent time start should be smaller than end")
	ErrUserNameEmpty            = errors.New("user name is empty")
	ErrUserNameTooLong          = errors.New("user name is too long")
	ErrUserEmailEmpty           = errors.New("user email is empty")
	ErrUserEmailFormat          = errors.New("user email format is invalid")
	ErrUserPhoneLen             = errors.New("user phone length should be 10 or 11")
	ErrUserPhoneFormat          = errors.New("user phone should contains 10 or 11 numbers")
	ErrRulePatternEmpty         = errors.New("rule pattern is empty")
	ErrRulePatternTooLong       = errors.New("rule pattern is too long")
	ErrRulePatternContainsSpace = errors.New("rule pattern contains spaces")
	ErrRulePatternFormat        = errors.New("rule pattern format is invalid")
	ErrMetricNameEmpty          = errors.New("metric name is empty")
	ErrMetricNameTooLong        = errors.New("metric name is too long")
	ErrMetricStampTooSmall      = errors.New("metric stamp is too small")
)

// ValidateProjectName validates project name
func ValidateProjectName(name string) error {
	if len(name) == 0 {
		// Empty
		return ErrProjectNameEmpty
	}
	if len(name) > MaxProjectNameLen {
		// Too long
		return ErrProjectNameTooLong
	}
	return nil
}

// ValidateProjectSilentRange validates project silent time start and end.
func ValidateProjectSilentRange(start, end int) error {
	if start < 0 || start > 23 {
		// Invalid number.
		return ErrProjectSilentTimeStart
	}
	if end < 0 || end > 23 {
		// Invalid number.
		return ErrProjectSilentTimeEnd
	}
	if start >= end {
		// Start is larger than end.
		return ErrProjectSilentTimeRange
	}
	return nil
}

// ValidateUserName validates user name.
func ValidateUserName(name string) error {
	if len(name) == 0 {
		// Empty
		return ErrUserNameEmpty
	}
	if len(name) > MaxUserNameLen {
		// Too long
		return ErrUserNameTooLong
	}
	return nil
}

// ValidateUserEmail validates user email.
func ValidateUserEmail(email string) error {
	if len(email) == 0 {
		// Empty
		return ErrUserEmailEmpty
	}
	if !strings.Contains(email, "@") {
		// Invalid format.
		return ErrUserEmailFormat
	}
	return nil
}

// ValidateUserPhone validates user phone.
func ValidateUserPhone(phone string) error {
	if len(phone) != 10 && len(phone) != 11 {
		// Invalid length.
		return ErrUserPhoneLen
	}
	if ok, _ := regexp.MatchString("^\\d{10,11}", phone); !ok {
		// Invalid format.
		return ErrUserPhoneFormat
	}
	return nil
}

// ValidateRulePattern validates rule pattern.
func ValidateRulePattern(pattern string) error {
	if len(pattern) == 0 {
		// Empty
		return ErrRulePatternEmpty
	}
	if len(pattern) > MaxRulePatternLen {
		// Too long
		return ErrRulePatternTooLong
	}
	if strings.ContainsAny(pattern, " \t\r\n") {
		// Contains space
		return ErrRulePatternContainsSpace
	}
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '*' {
			if i != 0 && pattern[i-1] != '.' {
				// Invalid format
				return ErrRulePatternFormat
			}
			if i != len(pattern)-1 && pattern[i+1] != '.' {
				// Invalid format
				return ErrRulePatternFormat
			}
		}
	}
	return nil
}

// ValidateMetricName validates metric name.
func ValidateMetricName(name string) error {
	if len(name) == 0 {
		// Empty
		return ErrMetricNameEmpty
	}
	if len(name) > MaxMetricNameLen {
		// Too long.
		return ErrMetricNameTooLong
	}
	return nil
}

// ValidateMetricStamp validates metric stamp.
func ValidateMetricStamp(stamp uint32) error {
	if stamp < MinMetricStamp {
		// Too small.
		return ErrMetricStampTooSmall
	}
	return nil
}
