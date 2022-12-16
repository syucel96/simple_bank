package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syucel96/simplebank/util"
)

func createRandomTransfer(t *testing.T, fromAccountId, toAccountId *int64) Transfer {
	arg := CreateTransferParams{
		FromAccountID: *fromAccountId,
		ToAccountID:   *toAccountId,
		Amount:        util.RandomMoney(true),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	account := createRandomAccount(t, user1.Username)
	account2 := createRandomAccount(t, user2.Username)
	createRandomTransfer(t, &account.ID, &account2.ID)
}

func TestGetTransfer(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)
	account := createRandomAccount(t, user1.Username)
	account2 := createRandomAccount(t, user2.Username)
	createdTransfer := createRandomTransfer(t, &account.ID, &account2.ID)

	transfer, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, createdTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, createdTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, createdTransfer.Amount, transfer.Amount)
	require.Equal(t, createdTransfer.ID, transfer.ID)
	require.WithinDuration(t, createdTransfer.CreatedAt, transfer.CreatedAt, time.Microsecond)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 5; i++ {
		user1 := createRandomUser(t)
		user2 := createRandomUser(t)
		account := createRandomAccount(t, user1.Username)
		account2 := createRandomAccount(t, user2.Username)
		createRandomTransfer(t, &account.ID, &account2.ID)
		createRandomTransfer(t, &account2.ID, &account.ID)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
