// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/util/assert"
	"math/rand"
	"testing"
)

func genLongString(length int) string {
	d := "abcdefghijk0123456789"
	s := ""
	for i := 0; i < length; i++ {
		s = s + string(d[rand.Intn(len(d))])
	}
	return s
}

func TestValidateProjectName(t *testing.T) {
	assert.Ok(t, ValidateProjectName("") == ErrProjectNameEmpty)
	assert.Ok(t, ValidateProjectName(genLongString(MaxProjectNameLen+1)) == ErrProjectNameTooLong)
	assert.Ok(t, ValidateProjectName("project") == nil)
}

func TestValidateProjectSilentTimeRange(t *testing.T) {
	assert.Ok(t, ValidateProjectSilentRange(39, 6) == ErrProjectSilentTimeStart)
	assert.Ok(t, ValidateProjectSilentRange(0, 29) == ErrProjectSilentTimeEnd)
	assert.Ok(t, ValidateProjectSilentRange(7, 4) == ErrProjectSilentTimeRange)
	assert.Ok(t, ValidateProjectSilentRange(1, 9) == nil)
}

func TestValidateUserName(t *testing.T) {
	assert.Ok(t, ValidateUserName("") == ErrUserNameEmpty)
	assert.Ok(t, ValidateUserName(genLongString(MaxUserNameLen+1)) == ErrUserNameTooLong)
	assert.Ok(t, ValidateUserName("user") == nil)
}

func TestValidateUserEmail(t *testing.T) {
	assert.Ok(t, ValidateUserEmail("") == ErrUserEmailEmpty)
	assert.Ok(t, ValidateUserEmail("abc") == ErrUserEmailFormat)
	assert.Ok(t, ValidateUserEmail("hit9@ele.me") == nil)
}

func TestValidateUserPhone(t *testing.T) {
	assert.Ok(t, ValidateUserPhone("123456789012") == ErrUserPhoneLen)
	assert.Ok(t, ValidateUserPhone("12345678a01") == ErrUserPhoneFormat)
	assert.Ok(t, ValidateUserPhone("18701616177") == nil)
}

func TestValidateRulePattern(t *testing.T) {
	assert.Ok(t, ValidateRulePattern("") == ErrRulePatternEmpty)
	assert.Ok(t, ValidateRulePattern("abc efg") == ErrRulePatternContainsSpace)
	assert.Ok(t, ValidateRulePattern("abc*.s") == ErrRulePatternFormat)
	assert.Ok(t, ValidateRulePattern("abc.*.s") == nil)
	assert.Ok(t, ValidateRulePattern("abc.*.*") == nil)
	assert.Ok(t, ValidateRulePattern("*.abc.*") == nil)
}

func TestValidateMetricName(t *testing.T) {
	assert.Ok(t, ValidateMetricName("") == ErrMetricNameEmpty)
	assert.Ok(t, ValidateMetricName(genLongString(MaxMetricNameLen+1)) == ErrMetricNameTooLong)
}

func TestValidateMetricStamp(t *testing.T) {
	assert.Ok(t, ValidateMetricStamp(123) == ErrMetricStampTooSmall)
}
