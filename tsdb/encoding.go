package tsdb

import (
	"fmt"
	"strconv"
)

const (
	// Namespace
	stampKeyLeader = 'S'
	nameKeyLeader  = 'N'
	// Timestamp
	stampHorizon   uint64 = 1449308016 // further stamps must be larger than this
	stampConvBase         = 36         // convert to 36-hex string
	stampStringLen        = 7          // enough for 100 years
)

var (
	ErrNameKey    = NewErrCorruptedWithString("invalid series name key")
	ErrNameValue  = NewErrCorruptedWithString("invalid series name value")
	ErrStampKey   = NewErrCorruptedWithString("invalid timestamp key")
	ErrStampValue = NewErrCorruptedWithString("invalid timestamp value")
)

// Encode series name to leveldb key.
func encodeNameKey(name string) string {
	return fmt.Sprintf("%c%s", nameKeyLeader, name)
}

// Decode name key to series name.
func decodeNameKey(key string) (string, error) {
	if len(key) < 2 {
		return "", ErrNameKey
	}
	return key[1:], nil
}

// Encode name value with trend and stamp.
func encodeNameValue(trend float64, stamp uint64) string {
	return fmt.Sprintf("%.3f:%d", trend, stamp)
}

// Decode name value to trend and stamp.
func decodeNameValue(s string) (trend float64, stamp uint64, err error) {
	n, err := fmt.Sscanf(s, "%f:%d", &trend, &stamp)
	if err != nil {
		err = NewErrCorrupted(err)
		return
	}
	if n != 2 {
		err = ErrNameValue
		return
	}
	return
}

// Encode stamp to leveldb key with name.
func encodeStampKey(name string, stamp uint64) string {
	s := strconv.FormatUint(stamp-stampHorizon, stampConvBase)
	return fmt.Sprintf("%c%s%0*s", stampKeyLeader, name, stampStringLen, s)
}

// Decode stamp key to name and stamp.
func decodeStampKey(key string) (name string, stamp uint64, err error) {
	if len(key) <= stampStringLen {
		err = ErrStampKey
		return
	}
	i := len(key) - stampStringLen
	name = key[1:i]
	// Assume to be a uint64
	n, err := strconv.ParseUint(key[i:], stampConvBase, 64)
	if err != nil {
		err = NewErrCorrupted(err)
		return
	}
	stamp = n + stampHorizon
	return
}

// Encode stamp value to string with value and score.
func encodeStampValue(value float64, score float64) string {
	return fmt.Sprintf("%.3f:%.3f", value, score)
}

// Decode stamp value string to value and score.
func decodeStampValue(s string) (value, score float64, err error) {
	n, err := fmt.Sscanf(s, "%f:%f", &value, &score)
	if err != nil {
		err = NewErrCorrupted(err)
		return
	}
	if n != 2 {
		err = ErrStampValue
		return
	}
	return
}
