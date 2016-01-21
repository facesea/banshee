// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"regexp"
	"strings"
)

// Limitations
const maxProjectNameLen = 64
const maxUserNameLen = 32
const maxRulePatternLen = 256

// validate project name.
func validateProjectName(name string) *WebError {
	if len(name) <= 0 {
		// Empty
		return ErrProjectNameEmpty
	}
	if len(name) > maxProjectNameLen {
		// Too long
		return ErrProjectNameTooLong
	}
	return nil
}

// validate user name.
func validateUserName(name string) *WebError {
	if len(name) <= 0 {
		// Empty
		return ErrUserNameEmpty
	}
	if len(name) > maxUserNameLen {
		// Too long
		return ErrUserNameTooLong
	}
	return nil
}

// validate user email.
func validateUserEmail(email string) *WebError {
	if len(email) <= 0 {
		// Empty
		return ErrUserEmailEmpty
	}
	if !strings.Contains(email, "@") {
		// Invalid
		return ErrUserEmail
	}
	return nil
}

// validate user phone.
func validateUserPhone(phone string) *WebError {
	if len(phone) != 10 && len(phone) != 11 {
		// Bad length.
		return ErrUserPhoneLen
	}
	if ok, _ := regexp.MatchString("^\\d{10,11}", phone); !ok {
		// Bad format.
		return ErrUserPhone
	}
	return nil
}

// validate rule pattern.
func validateRulePattern(pattern string) *WebError {
	if len(pattern) <= 0 {
		// Empty
		return ErrRulePatternEmpty
	}
	if len(pattern) > maxRulePatternLen {
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
				// Bad format.
				return ErrRulePattern
			}
			if i != len(pattern)-1 && pattern[i+1] != '.' {
				// Bad format.
				return ErrRulePattern
			}
		}
	}
	return nil
}
