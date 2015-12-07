package tsdb

import "testing"

func TestIsCorrupted(t *testing.T) {
	if IsErrCorrupted(ErrNotFound) {
		t.Errorf("IsErrCorrupted(ErrNotFound) should return false")
	}
	if !IsErrCorrupted(NewErrCorruptedWithString("something wrong")) {
		t.Errorf("IsErrCorrupted(NewErrCorruptedWithString(s)) should return true")
	}
}
