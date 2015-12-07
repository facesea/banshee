package tsdb

import "errors"

var (
	// Returned if a key is in invalid format
	ErrKeyFormat = errors.New("invalid key format")
)
