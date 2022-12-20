package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syucel96/simplebank/util"
)

func createRandomAccount(t *testing.T, username string, fields ...string) Account {
	var arg CreateAccountParams
	if len(fields) == 0 {
		arg = CreateAccountParams{
			Owner:    username,
			Balance:  util.RandomMoney(false),
			Currency: util.RandomCurrency(),
		}
	} else {
		arg = CreateAccountParams{
			Owner:    username,
			Balance:  util.RandomMoney(false),
			Currency: fields[0],
		}
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
	user := createRandomUser(t)
	currencies := []string{util.USD, util.EUR, util.GBP, util.CAD, util.JPY, util.TRY}
	for i := range currencies {
		createRandomAccount(t, user.Username, currencies[i])
	}

	limit := int32(util.RandomInt(int64(2), int64(6)))
	offset := int32(util.RandomInt(int64(0), int64(3)))

	var expected int32
	if limit+offset > 6 {
		expected = 6 - offset
	} else {
		expected = limit
	}

	arg := ListAccountsParams{
		Owner:  user.Username,
		Limit:  limit,
		Offset: offset,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, int(expected))

	for i, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, currencies[i+int(offset)], account.Currency)
		require.Equal(t, user.Username, account.Owner)
	}
}
