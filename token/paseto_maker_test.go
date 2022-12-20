package token

import (
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/stretchr/testify/require"
	"github.com/syucel96/simplebank/util"
)

func TestNewPasetoMaker(t *testing.T) {
	invTokenMin, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize - 1))
	require.Error(t, err)
	require.Empty(t, invTokenMin)

	invTokenMax, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize + 1))
	require.Error(t, err)
	require.Empty(t, invTokenMax)

	maker, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)
}
