package tsdb

import "errors"

var (
	// ErrKeyFormat is returned if the key is invalid when being decode.
	ErrKeyFormat = errors.New("Invalid db key format")
)
