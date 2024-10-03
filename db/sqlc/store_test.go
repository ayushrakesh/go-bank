package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createTestAccount(t)
	account2 := createTestAccount(t)
	amount := int64(10)
	n := 5

	resCh := make(chan TransferTxResult)
	errCh := make(chan error)
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.WithValue(context.Background(), txName, i+1)
			res, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errCh <- err
			resCh <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errCh
		require.NoError(t, err)

		result := <-resCh
		require.NotEmpty(t, result)

		//checking transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//checking from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//checking to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// checkupdated accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	checkAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	checkAccount2, errr := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, errr)

	require.Equal(t, account1.Balance-int64(n)*amount, checkAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, checkAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createTestAccount(t)
	account2 := createTestAccount(t)
	amount := int64(10)
	n := 10

	// resCh := make(chan TransferTxResult)
	errCh := make(chan error)
	// existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errCh <- err
			// resCh <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errCh
		require.NoError(t, err)
	}

	checkAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	checkAccount2, errr := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, errr)

	require.Equal(t, account1.Balance, checkAccount1.Balance)
	require.Equal(t, account2.Balance, checkAccount2.Balance)
}
