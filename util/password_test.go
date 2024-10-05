package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPass1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass1)

	err = CheckPassword(password, hashedPass1)
	require.NoError(t, err)

	wrongPass := RandomString(6)
	err = CheckPassword(wrongPass, hashedPass1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPass2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPass2)
	require.NotEqual(t, hashedPass1, hashedPass2)
}
