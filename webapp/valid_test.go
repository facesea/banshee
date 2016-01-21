// Copyright 2015 Eleme Inc. All rights reserved.

package webapp

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
	assert.Ok(t, validateProjectName("") == ErrProjectNameEmpty)
	assert.Ok(t, validateProjectName(genLongString(maxProjectNameLen+1)) == ErrProjectNameTooLong)
	assert.Ok(t, validateProjectName("project") == nil)
}

func TestValidateUserName(t *testing.T) {
	assert.Ok(t, validateUserName("") == ErrUserNameEmpty)
	assert.Ok(t, validateUserName(genLongString(maxUserNameLen+1)) == ErrUserNameTooLong)
	assert.Ok(t, validateUserName("user") == nil)
}

func TestValidateUserEmail(t *testing.T) {
	assert.Ok(t, validateUserEmail("") == ErrUserEmailEmpty)
	assert.Ok(t, validateUserEmail("abc") == ErrUserEmail)
	assert.Ok(t, validateUserEmail("hit9@ele.me") == nil)
}

func TestValidateUserPhone(t *testing.T) {
	assert.Ok(t, validateUserPhone("123456789012") == ErrUserPhoneLen)
	assert.Ok(t, validateUserPhone("12345678a01") == ErrUserPhone)
	assert.Ok(t, validateUserPhone("18701616177") == nil)
}

func TestValidateRulePattern(t *testing.T) {
	assert.Ok(t, validateRulePattern("") == ErrRulePatternEmpty)
	assert.Ok(t, validateRulePattern("abc efg") == ErrRulePatternContainsSpace)
	assert.Ok(t, validateRulePattern("abc*.s") == ErrRulePattern)
	assert.Ok(t, validateRulePattern("abc.*.s") == nil)
	assert.Ok(t, validateRulePattern("abc.*.*") == nil)
	assert.Ok(t, validateRulePattern("*.abc.*") == nil)
}
