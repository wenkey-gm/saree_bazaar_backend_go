package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse")
	require.NoError(t, err)

	ok, err := ComparePasswords(hash, "correct horse")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = ComparePasswords(hash, "wrong password")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestComparePasswordsRejectsMalformedHash(t *testing.T) {
	ok, err := ComparePasswords("no-salt-here", "anything")
	assert.Error(t, err)
	assert.False(t, ok)
}
