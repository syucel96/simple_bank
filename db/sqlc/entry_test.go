package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syucel96/simplebank/util"
)

func createRandomEntry(t *testing.T, accountId *int64) Entry {
	arg := CreateEntryParams{
		AccountID: *accountId,
		Amount:    util.RandomMoney(true),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, &account.ID)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	createdEntry := createRandomEntry(t, &account.ID)

	entry, err := testQueries.GetEntry(context.Background(), createdEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, createdEntry.AccountID, entry.AccountID)
	require.Equal(t, createdEntry.Amount, entry.Amount)
	require.Equal(t, createdEntry.ID, entry.ID)
	require.WithinDuration(t, createdEntry.CreatedAt, entry.CreatedAt, time.Microsecond)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		account := createRandomAccount(t)
		createRandomEntry(t, &account.ID)
		createRandomEntry(t, &account.ID)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
