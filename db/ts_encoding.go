package tsdb

import (
	"fmt"
	"strconv"
)

const (
	// Further incoming timestamp will be stored as the diff to this horizon.
	timeStampHorizon uint64 = 1449308016
	// Timestamps will be converted to 36-hex string format before the are put
	// into db,
	timeStampConvBase = 36
	// A timestamp in 36-hex format string with this length is enough for the
	// further 100 years to use.
	timeStampSLength = 7
)

var (
	ErrTimeStampString = NewErrCorruptedWithString("invalid timestamp string")
	ErrTsKey           = NewErrCorruptedWithString("invalid series key")
	ErrTsName          = NewErrCorruptedWithString("invalid series name")
)

// Encode timestamp to 36-hex string format.
func encodeTimeStamp(t uint64) string {
	diff := t - timeStampHorizon
	return strconv.FormatUint(diff, timeStampConvBase)
}

// Decode timestamp from a 36-hex string.
func decodeTimeStamp(s string) (uint64, error) {
	if len(s) != timeStampSLength {
		return 0, ErrTimeStampString
	}
	// Assume to be a uint64 integer
	n, err := strconv.ParseUint(s, timeStampConvBase, 64)
	if err != nil {
		return 0, NewErrCorrupted(err)
	}
	return n + timeStampHorizon, nil
}

// Encode timeseries name and timestamp to db key.
func encodeTsKey(name string, t uint64) string {
	s := encodeTimeStamp(t)
	return fmt.Sprintf("%c%s%s", nsTsKey, name, s)
}

// Decode timeseries data key into name and timestamp.
func decodeTsKey(s string) (string, uint64, error) {
	if len(s) <= timeStampSLength {
		return "", 0, ErrTsKey
	}
	idx := len(s) - timeStampSLength
	name := s[1:idx]
	t, err := decodeTimeStamp(s[idx:])
	if err != nil {
		return "", 0, NewErrCorrupted(err)
	}
	return name, t, nil
}

// Encode timeseries name to db key.
func encodeTsName(name string) string {
	return fmt.Sprintf("%c%s", nsTsName, name)
}

// Decode timeseries name from db key.
func decodeTsName(s string) (string, error) {
	if len(s) < 2 {
		return "", ErrTsName
	}
	return s[1:], nil
}

// Encode timeseries value to db value.
func encodeTsValue(v float64) string {
	return fmt.Sprintf("%.3f", v)
}

// Decode timeseries value from db value.
func decodeTsValue(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v, NewErrCorrupted(err)
	}
	return v, nil
}
