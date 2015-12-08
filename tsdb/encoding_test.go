package tsdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodingPrefixes(t *testing.T) {
	s := fmt.Sprintf("%c%c%c%c", prefixTsName, prefixTsKey, prefixHashName, prefixHashKey)
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			assert.NotEqual(t, s[i], s[j])
		}
	}
}

func TestEncodeTimeStamp(t *testing.T) {
	timeStamp := uint64(1449481973)
	excepted := "0003q85"
	actually := encodeTimeStamp(timeStamp)
	assert.Equal(t, actually, excepted)
}

func TestDecodeTimeStamp(t *testing.T) {
	excepted := uint64(1449481973)
	actually, err := decodeTimeStamp("0003q85")
	assert.Nil(t, err)
	assert.Equal(t, actually, excepted)
}

func TestEnDecodeTimeStamp(t *testing.T) {
	timeStamp := uint64(1449483044)
	s := encodeTimeStamp(timeStamp)
	assert.True(t, len(s) == 7)
	actually, err := decodeTimeStamp(s)
	excepted := timeStamp
	assert.Nil(t, err)
	assert.Equal(t, actually, excepted)
}

func TestEncodeTsKey(t *testing.T) {
	excepted := fmt.Sprintf("%c%s%s", prefixTsKey, "name", "0003q85")
	actually := encodeTsKey("name", uint64(1449481973))
	assert.Equal(t, actually, excepted)
}

func TestDecodeTsKey(t *testing.T) {
	exceptedName := "name"
	exceptedTimeStamp := uint64(1449481973)
	key := fmt.Sprintf("%c%s%s", prefixTsKey, "name", "0003q85")
	actuallyName, actuallyTimeStamp, err := decodeTsKey(key)
	assert.Nil(t, err)
	assert.Equal(t, actuallyName, exceptedName)
	assert.Equal(t, actuallyTimeStamp, exceptedTimeStamp)
}

func TestEnDecodeTsKey(t *testing.T) {
	name := "ts"
	timeStamp := uint64(1449484137)
	key := encodeTsKey(name, timeStamp)
	assert.True(t, len(key) == 1+len(name)+timeStampSLength)
	exceptedName := name
	exceptedTimeStamp := timeStamp
	actuallyName, actuallyTimeStamp, err := decodeTsKey(key)
	assert.Nil(t, err)
	assert.Equal(t, actuallyName, exceptedName)
	assert.Equal(t, actuallyTimeStamp, exceptedTimeStamp)
}

func TestEncodeTsName(t *testing.T) {
	name := "testanc"
	excepted := fmt.Sprintf("%c%s", prefixTsName, name)
	actually := encodeTsName(name)
	assert.Equal(t, actually, excepted)
}

func TestDecodeTsName(t *testing.T) {
	name := "abctest"
	excepted := name
	actually, err := decodeTsName(encodeTsName(name))
	assert.Nil(t, err)
	assert.Equal(t, actually, excepted)
}

func TestEncodeValue(t *testing.T) {
	value := 1.28912
	excepted := "1.289"
	actually := encodeTsValue(value)
	assert.Equal(t, actually, excepted)
}

func TestDecodeTsValue(t *testing.T) {
	s := "1.389"
	excepted := 1.389
	actually, err := decodeTsValue(s)
	assert.Nil(t, err)
	assert.Equal(t, actually, excepted)
}
