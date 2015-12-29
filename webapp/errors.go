// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

import (
	"errors"
	"fmt"
	"net/http"
)

// WebError is errors for web operations.
type WebError struct {
	// HTTP status code
	code int
	// Error
	err error
}

// Errors
var (
	// Common
	ErrBadRequest = NewWebErrorWithText(http.StatusBadRequest, "Bad request")
	ErrNotNull    = NewWebErrorWithText(http.StatusBadRequest, "Null value")
	ErrPrimaryKey = NewWebErrorWithText(http.StatusForbidden, "Primarykey voilated")
	ErrUnique     = NewWebErrorWithText(http.StatusForbidden, "Value should be unique")
	ErrNotFound   = NewWebErrorWithText(http.StatusNotFound, "Not found")
	// Project
	ErrProjectID            = NewWebErrorWithText(http.StatusBadRequest, "Bad project id")
	ErrProjectName          = NewWebErrorWithText(http.StatusBadRequest, "Bad project name")
	ErrProjectNotFound      = NewWebErrorWithText(http.StatusNotFound, "Project not found")
	ErrDuplicateProjectName = NewWebErrorWithText(http.StatusForbidden, "Duplicate project name")
	// User
	ErrUserID            = NewWebErrorWithText(http.StatusBadRequest, "Bad user id")
	ErrUserName          = NewWebErrorWithText(http.StatusBadRequest, "Bad user name")
	ErrUserEmail         = NewWebErrorWithText(http.StatusBadRequest, "Bad user email")
	ErrUserPhone         = NewWebErrorWithText(http.StatusBadRequest, "Bad user phone")
	ErrUserNotFound      = NewWebErrorWithText(http.StatusNotFound, "User not found")
	ErrDuplicateUserName = NewWebErrorWithText(http.StatusForbidden, "Duplicate user name")
	// Rule
	ErrRuleID               = NewWebErrorWithText(http.StatusBadRequest, "Bad rule id")
	ErrRulePattern          = NewWebErrorWithText(http.StatusBadRequest, "Bad rule pattern")
	ErrRuleWhen             = NewWebErrorWithText(http.StatusBadRequest, "Bad rule condition")
	ErrRuleProjectID        = NewWebErrorWithText(http.StatusBadRequest, "Bad rule project id")
	ErrDuplicateRulePattern = NewWebErrorWithText(http.StatusForbidden, "Duplicate rule pattern")
	ErrRuleNotFound         = NewWebErrorWithText(http.StatusNotFound, "Rule not found")
	// Metric
	ErrMetricNotFound = NewWebErrorWithText(http.StatusNotFound, "Metric not found")
)

// NewWebError creates a WebError.
func NewWebError(code int, err error) *WebError {
	return &WebError{code, err}
}

// NewWebErrorWithText creates a WebError with text.
func NewWebErrorWithText(code int, text string) *WebError {
	return &WebError{code, errors.New(text)}
}

// Error returns the string format of the WebError.
func (err *WebError) Error() string {
	return fmt.Sprintf("[%d]: %s", err.code, err.err.Error())
}

// NewUnexceptedWebError returns an unexcepted WebError.
func NewUnexceptedWebError(err error) *WebError {
	return NewWebError(http.StatusInternalServerError, err)
}
