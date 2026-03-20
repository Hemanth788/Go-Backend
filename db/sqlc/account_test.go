package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"go.com/go-backend/util"
)

func createRandomAccount(t *testing.T) (CreateAccountParams, Account) {
	accountArgs := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomMoney(),
	}

	account, err := testQueries.CreateAccount(context.Background(), accountArgs)
	require.NoError(t, err)

	return accountArgs, account
}

func TestCreateAccount(t *testing.T) {
	accountArgs, account := createRandomAccount(t)

	require.NotEmpty(t, account)

	require.Equal(t, accountArgs.Owner, account.Owner)
	require.Equal(t, accountArgs.Balance, account.Balance)
	require.Equal(t, accountArgs.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	_, account1 := createRandomAccount(t)
	accountFetched, errFetched := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, errFetched)
	require.NotEmpty(t, accountFetched)

	require.EqualValues(t, account1, accountFetched)
}

func TestUpdateAccount(t *testing.T) {
	_, account1 := createRandomAccount(t)

	updateArgs := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	accountUpdated, err := testQueries.UpdateAccount(context.Background(), updateArgs)

	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, account1.ID, accountUpdated.ID)
	require.Equal(t, account1.Owner, accountUpdated.Owner)
	require.Equal(t, account1.Currency, accountUpdated.Currency)
	require.Equal(t, updateArgs.Balance, accountUpdated.Balance)
}

func TestDeleteAccount(t *testing.T) {
	_, account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	accountDeleted, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountDeleted)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
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
