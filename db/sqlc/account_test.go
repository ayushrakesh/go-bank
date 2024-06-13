package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ayushrakesh/go-bank/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	ar := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), ar)
	require.NotEmpty(t, account)
	require.NoError(t, err)

	require.Equal(t, ar.Balance, account.Balance)
	require.Equal(t, ar.Currency, account.Currency)
	require.Equal(t, ar.Owner, account.Owner)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createTestAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)

	require.NotZero(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)

	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createTestAccount(t)

	ar := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}
	acc2, err := testQueries.UpdateAccount(context.Background(), ar)

	require.NoError(t, err)

	require.NotZero(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

	require.Equal(t, ar.Balance, acc2.Balance)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)

	require.NoError(t, err)

	acc2, er := testQueries.GetAccount(context.Background(), acc1.ID)

	require.Error(t, er)
	require.Empty(t, acc2)

	require.EqualError(t, er, sql.ErrNoRows.Error())

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}
	ar := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accs, err := testQueries.ListAccounts(context.Background(), ar)

	require.NoError(t, err)
	require.NotEmpty(t, accs)
	require.Len(t, accs, 5)

	for i := 0; i < len(accs); i++ {
		require.NotEmpty(t, accs[i])
	}
}
