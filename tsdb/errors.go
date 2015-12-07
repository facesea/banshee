package tsdb

import "errors"

// Following errors are user-level and should be checked, and other errors from
// tsdb are corruption errors. We can use IsErrCorrupted to check if an error
// is a corruption error.
var (
	ErrNotFound = errors.New("tsdb: not found")
)

// ErrCorrupted is the type that wraps errors that indicate corruption in
// the database.
type ErrCorrupted struct {
	Err error
}

func (err *ErrCorrupted) Error() string {
	return err.Error()
}

// NewErrCorrupted creates new ErrCorrupted error.
func NewErrCorrupted(err error) error {
	return &ErrCorrupted{err}
}

// NewErrCorrupted creates new ErrCorrupted error.
func NewErrCorruptedWithString(s string) error {
	return NewErrCorrupted(errors.New(s))
}

// IsErrCorrupted returns a boolean indicating whether the error is indicating
// a corruption.
func IsErrCorrupted(err error) bool {
	_, ok := err.(*ErrCorrupted)
	return ok
}
