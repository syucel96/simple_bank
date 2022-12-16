package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syucel96/simplebank/util"
)

func createRandomAccount(t *testing.T, username string) Account {
	arg := CreateAccountParams{
		Owner:    username,
		Balance:  util.RandomMoney(false),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	createRandomAccount(t, user.Username)
}

func TestGetAccount(t *testing.T) {
	user := createRandomUser(t)
	createdAccount := createRandomAccount(t, user.Username)

	account, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createdAccount.Owner, account.Owner)
	require.Equal(t, createdAccount.Balance, account.Balance)
	require.Equal(t, createdAccount.Currency, account.Currency)
	require.Equal(t, createdAccount.ID, account.ID)
	require.WithinDuration(t, createdAccount.CreatedAt, account.CreatedAt, time.Microsecond)
}

func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)
	createdAccount := createRandomAccount(t, user.Username)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: util.RandomMoney(false),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createdAccount.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, createdAccount.Currency, account.Currency)
	require.Equal(t, createdAccount.ID, account.ID)
	require.WithinDuration(t, createdAccount.CreatedAt, account.CreatedAt, time.Microsecond)
}

func TestDeleteAccount(t *testing.T) {
	user := createRandomUser(t)
	createdAccount := createRandomAccount(t, user.Username)

	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		user := createRandomUser(t)
		createRandomAccount(t, user.Username)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
