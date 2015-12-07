package tsdb

import (
	"fmt"
	"strconv"
)

const (
	// Further incoming timestamp will be stored as the diff to this horizon.
	timeStampHorizon = 1449308016
	// Timestamps will be converted to 36-hex string format before the are put
	// into db,
	timeStampConvBase = 36
	// A timestamp in 36-hex format string with this length is enough for the
	// further 100 years to use.
	timeStampSLength = 7
)

// Encode timestamp to 36-hex string format.
func encodeTimeStamp(t TimeStamp) string {
	diff := t - timeStampHorizon
	return strconv.FormatUint(diff, timeStampConvBase)
}

// Decode timestamp from a 36-hex string.
func decodeTimeStamp(s string) (TimeStamp, error) {
	if len(s) != timeStampSLength {
		return 0, ErrKeyFormat
	}
	// Assume to be a uint64 integer
	n, err := strconv.ParseUint(s, timeStampConvBase, 64)
	if err != nil {
		return 0, err
	}
	return n + timeStampHorizon
}

// Encode timeseries name and timestamp to db key.
func encodeTsKey(name string, t TimeStamp) string {
	s := encodeTimeStamp(t)
	return fmt.Sprintf("%c%s%s", nsTsKey, name, s)
}

// Decode timeseries data key into name and timestamp.
func decodeTsKey(key string) (string, TimeStamp, error) {
	if len(s) <= timeStampSLength {
		return "", 0, ErrKeyFormat
	}
	idx := len(s) - timeStampSLength
	name := s[1:idx]
	t, err := decodeTimeStamp(s[idx:])
	if err != nil {
		return "", 0, errr
	}
	return name, t, nil
}

// Encode timeseries name to db key.
func encodeTsName(name string) string {
	return fmt.Sprintf("%c%s", nsTsName, name, s)
}

// Decode timeseries name from db key.
func decodeTsName(s string) (string, error) {
	if len(s) < 2 {
		return "", ErrKeyFormat
	}
	return s[1:], nil
}

// Encode timeseries value to db value.
func encodeTsValue(v Value) string {
	return fmt.Sprintf("%.3f", v)
}

// Decode timeseries value from db value.
func decodeTsValue(s string) (Value, error) {
	return strconv.ParseFloat(s, 64)
}
