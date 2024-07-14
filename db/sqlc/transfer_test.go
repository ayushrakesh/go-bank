package db

import (
	"context"
	"testing"
	"time"

	"github.com/ayushrakesh/go-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from Account, to Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.ID)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	from := createTestAccount(t)
	to := createTestAccount(t)

	createRandomTransfer(t, from, to)
}

func TestGetTransfer(t *testing.T) {
	from := createTestAccount(t)
	to := createTestAccount(t)

	trn := createRandomTransfer(t, from, to)

	transfer, err := testQueries.GetTransfer(context.Background(), trn.ID)

	require.NoError(t, err)
	require.NotZero(t, transfer)

	require.Equal(t, trn.ID, transfer.ID)
	require.Equal(t, trn.FromAccountID, transfer.FromAccountID)
	require.Equal(t, trn.ToAccountID, transfer.ToAccountID)
	require.Equal(t, trn.Amount, transfer.Amount)

	require.WithinDuration(t, trn.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	from := createTestAccount(t)
	to := createTestAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, from, to)
	}
	arg := ListTransfersParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Limit:         5,
		Offset:        0,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)

		require.Equal(t, transfer.FromAccountID, from.ID)
		require.Equal(t, transfer.ToAccountID, to.ID)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
	}
}
