package tsdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyLeaders(t *testing.T) {
	assert.NotEqual(t, nameKeyLeader, stampKeyLeader)
}

func TestEncodeNameKey(t *testing.T) {
	name := "abc"
	excepted := fmt.Sprintf("%c%s", nameKeyLeader, name)
	actually := encodeNameKey(name)
	assert.Equal(t, excepted, actually)
}

func TestDecodeNameKey(t *testing.T) {
	name := "abctest"
	excepted := name
	actually, err := decodeNameKey(encodeNameKey(name))
	assert.Nil(t, err)
	assert.Equal(t, excepted, actually)
}

func TestEncodeNameValue(t *testing.T) {
	trend := 1.23934
	stamp := uint64(1449481979)
	excepted := "1.239:1449481979"
	actually := encodeNameValue(trend, stamp)
	assert.Equal(t, excepted, actually)
}

func TestDecodeNameValue(t *testing.T) {
	exceptedTrend := 1.38
	exceptedStamp := uint64(1449481993)
	s := encodeNameValue(exceptedTrend, exceptedStamp)
	actuallyTrend, actuallyStamp, err := decodeNameValue(s)
	assert.Nil(t, err)
	assert.Equal(t, exceptedTrend, actuallyTrend)
	assert.Equal(t, exceptedStamp, actuallyStamp)
}

func TestEncodeStampKey(t *testing.T) {
	name := "abc"
	stamp := uint64(1449481973)
	excepted := fmt.Sprintf("%c%s%s", stampKeyLeader, name, "0003q85")
	actually := encodeStampKey(name, stamp)
	assert.Equal(t, excepted, actually)
}

func TestDecodeStampKey(t *testing.T) {
	exceptedName := "abc"
	exceptedStamp := uint64(1449483044)
	s := encodeStampKey(exceptedName, exceptedStamp)
	actuallyName, actuallyStamp, err := decodeStampKey(s)
	assert.Nil(t, err)
	assert.Equal(t, exceptedName, actuallyName)
	assert.Equal(t, exceptedStamp, actuallyStamp)
}

func TestEncodeStampValue(t *testing.T) {
	value := 12.289
	score := 1.2
	excepted := "12.289:1.200"
	actually := encodeStampValue(value, score)
	assert.Equal(t, excepted, actually)
}

func TestDecodeStampValue(t *testing.T) {
	exceptedValue := 12.389
	exceptedScore := 0.999
	s := encodeStampValue(exceptedValue, exceptedScore)
	actuallyValue, actuallyScore, err := decodeStampValue(s)
	assert.Nil(t, err)
	assert.Equal(t, exceptedValue, actuallyValue)
	assert.Equal(t, exceptedScore, actuallyScore)
}
