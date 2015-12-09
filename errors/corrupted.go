// Copyright 2015 Eleme Inc. All rights reserved.
package errors

import "errors"

// ErrCorrupted is the type that wraps errors that indicate corruption in
// the database.
type ErrCorrupted struct {
	Err error
}

func (err *ErrCorrupted) Error() string {
	return err.Error()
}

// NewCorrupted creates new ErrCorrupted error.
func NewErrCorrupted(err error) error {
	return &ErrCorrupted{err}
}

// NewCorrupted creates new ErrCorrupted error.
func NewErrCorruptedWithString(s string) error {
	return NewErrCorrupted(errors.New(s))
}

// IsCorrupted returns a boolean indicating whether the error is indicating
// a corruption.
func IsCorrupted(err error) bool {
	_, ok := err.(*ErrCorrupted)
	return ok
}
