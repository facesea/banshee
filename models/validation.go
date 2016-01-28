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
)

// Errors
var (
	ErrProjectNameEmpty         = errors.New("project name is empty")
	ErrProjectNameTooLong       = errors.New("project name is too long")
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
)

// Validate project name
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

// Validate user name.
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

// Validate user email.
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

// Validate user phone.
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

// Validate rule pattern.
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
