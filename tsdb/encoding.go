package tsdb

import (
	"fmt"
	"strconv"
)

func EncodeTimeStamp(t TimeStamp) string {
	return strconv.FormatUint(t-TimeStampBase, TimeStampConvBase)
}

func DecodeTimeStamp(s string) (TimeStamp, error) {
	if len(s) < TimeStampKeyLength {
		return 0, ErrKeyFormat
	}
	n, err := strconv.ParseUint(s, TimeStampConvBase, 64)
	if err != nil {
		return 0, err
	}
	return n + TimeStampBase, nil
}

func EncodeTimeStampKey(name string, t TimeStamp) string {
	return fmt.Sprintf("%c%s%s", TimeStampKeyPrefix, name, EncodeTimeStamp(t))
}

func DecodeTimeStampKey(s string) (string, TimeStamp, error) {
	if len(s) < TimeStampKeyLength+1 {
		return "", 0, ErrKeyFormat
	}
	name := s[1 : len(s)-TimeStampKeyLength]
	ts := s[len(s)-TimeStampKeyLength:]
	t, err := DecodeTimeStamp(ts)
	if err != nil {
		return "", 0, err
	}
	return name, t, nil
}

func EncodeSeriesNameKey(name string) string {
	return fmt.Sprintf("%c%s", SeriesNameKeyPrefix, name)
}

func DecodeSeriesNameKey(s string) (string, error) {
	if len(s) < 2 {
		return "", ErrKeyFormat
	}
	return s[1:], nil
}

func EncodeValue(v Value) string {
	return fmt.Sprintf("%.3f", v)
}

func DecodeValue(s string) (Value, error) {
	return strconv.ParseFloat(s, 64)
}
