package database

import (
	"context"
	"testing"

	"github.com/alansory/gobank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	return user
}
