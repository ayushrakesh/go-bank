package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

var txName = struct{}{}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		rber := tx.Rollback()
		if rber != nil {
			return fmt.Errorf("tx error %v, rb error %v", err, rber)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// txname := ctx.Value(txName)

	err := store.execTx(ctx, func(q *Queries) error {
		var errr error

		// fmt.Println(txname, "create transfer")
		result.Transfer, errr = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if errr != nil {
			return errr
		}

		// fmt.Println(txname, "create entry 1")
		result.FromEntry, errr = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if errr != nil {
			return errr
		}

		// fmt.Println(txname, "create entry 2")
		result.ToEntry, errr = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if errr != nil {
			return errr
		}

		// TODO update accounts
		// fmt.Println(txname, "get account 1")
		// fmt.Println(txname, "update account 1")
		result.FromAccount, errr = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if errr != nil {
			return errr
		}

		// // TODO update accounts
		// fmt.Println(txname, "get account 2")
		// fmt.Println(txname, "update account 2")
		result.ToAccount, errr = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if errr != nil {
			return errr
		}

		return nil
	})

	return result, err
}
