// Copyright 2015 Eleme Inc. All rights reserved.
package errors

import "errors"

// ErrFatal implements error, it indicates the error is fatal that progress
// should be interrupted.
type ErrFatal struct {
	Err error
}

func (err *ErrFatal) Error() string {
	return err.Err.Error()
}

// NewErrFatal creates an ErrFatal.
func NewErrFatal(err error) error {
	return &ErrFatal{err}
}

// NewErrFatalWithString creates new ErrFatal with string.
func NewErrFatalWithString(text string) error {
	return NewErrFatal(errors.New(text))
}

// IsFatal returns a boolean indicating whether the error is fatal.
func IsFatal(err error) bool {
	_, ok := err.(*ErrFatal)
	return ok
}
