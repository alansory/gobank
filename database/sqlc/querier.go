package database

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error)
	GetAccount(ctx context.Context, id int64) (Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error)
}

var _ Querier = (*Queries)(nil)
