package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.com/go-backend/util"
)

func createRandomUser(t *testing.T) (CreateUserParams, User) {
	userArgs := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "",
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err)

	return userArgs, user
}

func TestCreateUser(t *testing.T) {
	userArgs, user := createRandomUser(t)

	require.NotEmpty(t, user)

	require.Equal(t, userArgs.Username, user.Username)
	require.Equal(t, userArgs.Email, user.Email)
	require.Equal(t, userArgs.FullName, user.FullName)
	require.Equal(t, userArgs.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	_, user := createRandomUser(t)
	userFetched, errFetched := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, errFetched)
	require.NotEmpty(t, userFetched)

	require.EqualValues(t, user, userFetched)
}
