package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Why store? because Queries has CRUD for individual tables in the db but we need a place to do txns across tables and this is it
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, txFn func(*Queries) error) error {
	dbTxn, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// New takes any object that implements an interface DBTX which includes both sql.DB and sq.Tx
	q := New(dbTxn)
	err = txFn(q)

	if err != nil {
		if rbErr := dbTxn.Rollback(); rbErr != nil {
			return fmt.Errorf("dbTxn err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return dbTxn.Commit()
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

var TX_KEY = struct{}{}

func updateToAccount(ctx context.Context, result *TransferTxResult, q *Queries, arg TransferTxParams) error {
	var err error
	result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     arg.ToAccountID,
		Amount: +arg.Amount,
	})

	return err
}

func updateFromAccount(ctx context.Context, result *TransferTxResult, q *Queries, arg TransferTxParams) error {
	var err error
	result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     arg.FromAccountID,
		Amount: -arg.Amount,
	})

	return err
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(TX_KEY)

		fmt.Println("create transfer: ", txName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println("create from entry: ", txName)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println("create to entry: ", txName)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println("get from account for update: ", txName)
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

		// if err != nil {
		// 	return err
		// }

		if arg.FromAccountID < arg.ToAccountID {
			// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			// 	ID: arg.ToAccountID,
			// 	Balance: account1.Balance - arg.Amount,
			// })
			fmt.Println("update from account: ", txName)
			err = updateFromAccount(ctx, &result, q, arg)

			if err != nil {
				return err
			}

			// fmt.Println("get to account for update: ", txName)
			// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)

			// if err != nil {
			// 	return err
			// }
			fmt.Println("update to account: ", txName)
			err = updateToAccount(ctx, &result, q, arg)

			if err != nil {
				return err
			}
		} else {
			fmt.Println("update to account: ", txName)
			err = updateToAccount(ctx, &result, q, arg)

			fmt.Println("update from account: ", txName)
			err = updateFromAccount(ctx, &result, q, arg)
		}

		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID: arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })

		return nil
	})

	return result, err
}
