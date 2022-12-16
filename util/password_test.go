package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestEncryptAndDecryptPassword(t *testing.T) {
	password := RandomPassword()
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	require.NoError(t, CheckPassword(password, hashedPassword))

	wrongPassword := RandomPassword()
	require.EqualError(t, CheckPassword(wrongPassword, hashedPassword), bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)
}
