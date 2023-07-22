package database

import (
	"testing"

	"github.com/alansory/gobank/util"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := util.RandomString(6)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}
