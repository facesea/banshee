package tsdb

import (
	"fmt"
	"strconv"
)

const (
	// The prefix of the key to store time series name in db
	seriesNameKeyPrefix = '0'
	// The prefix of the key to store time stamp in db
	timeStampKeyPrefix = '1'
)

// The base time stamp, further incoming time stamps will be stored
// in db as the diff to this base.
const timeStampBase TimeStamp = 1449308016

// Time stamp will be converted to string in 36 hex format before put
// into the db.
const timeStampConvBase = 36

// This length of a 36-hex string is enough for the further 100 years
const timeStampKeyLength = 7

// Encode time stamp to string format in 36 hex format.
func encodeTimeStamp(t TimeStamp) string {
	return strconv.FormatUint(t-timeStampBase, timeStampConvBase)
}

// Decode time stamp from string. If the string is too short or
// invalid to parse into an unsigned integer, ErrKeyFormat will
// be returned.
func decodeTimeStamp(s string) (TimeStamp, error) {
	if len(s) < timeStampKeyLength {
		return 0, ErrKeyFormat
	}
	// Assume to be an uint64 integer
	n, err := strconv.ParseUint(s, timeStampConvBase, 64)
	if err != nil {
		return 0, err
	}
	return n + timeStampBase, nil
}

// Encode series name and time stamp to db key. A datapoint with
// time stamp and value will be stored in db in the following way:
//	<prefix><series-name><time stamp>	<value>
func encodeTimeStampKey(name string, t TimeStamp) string {
	s := encodeTimeStamp(t)
	return fmt.Sprintf("%c%s%s", timeStampKeyPrefix, name, s)
}

// Decode series name and time stamp from db key. If the key is too
// short or in a invalid format, ErrKeyFormat will be returned as an
// error.
func decodeTimeStampKey(s string) (string, TimeStamp, error) {
	if len(s) < timeStampKeyLength+1 {
		return "", 0, ErrKeyFormat
	}
	name := s[1 : len(s)-timeStampKeyLength]
	ts := s[len(s)-timeStampKeyLength:]
	t, err := decodeTimeStamp(ts)
	if err != nil {
		return "", 0, err
	}
	return name, t, nil
}

// Encode series name to db key. All series will be stored in db
// in the following way:
//	<prefix><series-name>	<value>
func encodeSeriesNameKey(name string) string {
	return fmt.Sprintf("%c%s", seriesNameKeyPrefix, name)
}

// Decode name from db key. If the string is too short, ErrKeyFormat
// will be returned as an error.
func decodeSeriesNameKey(s string) (string, error) {
	if len(s) < 2 {
		return "", ErrKeyFormat
	}
	return s[1:], nil
}

// Encode float64 typed value to string format.
func encodeValue(v Value) string {
	return fmt.Sprintf("%.3f", v)
}

// Decode series value from string.
func decodeValue(s string) (Value, error) {
	return strconv.ParseFloat(s, 64)
}
