package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syucel96/simplebank/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomPassword())
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomFullName(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createdUser.Username, user.Username)
	require.Equal(t, createdUser.HashedPassword, user.HashedPassword)
	require.Equal(t, createdUser.FullName, user.FullName)
	require.Equal(t, createdUser.Email, user.Email)
	require.WithinDuration(t, createdUser.PasswordChangedAt, user.PasswordChangedAt, time.Microsecond)
	require.WithinDuration(t, createdUser.CreatedAt, user.CreatedAt, time.Microsecond)
}

/* func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
} */
