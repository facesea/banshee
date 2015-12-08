package tsdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCorrupted(t *testing.T) {
	assert.False(t, IsErrCorrupted(ErrNotFound))
	assert.True(t, IsErrCorrupted(NewErrCorruptedWithString("something wrong")))
}
