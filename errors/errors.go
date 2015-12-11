// Copyright 2015 Eleme Inc. All rights reserved.

package errors

import "errors"

// New returns an normal error from text.
func New(text string) error {
	return errors.New(text)
}
