// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package database

import (
	"time"
)

type Account struct {
	ID        int64
	UserID    int64
	Balance   int64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Entry struct {
	ID        int64
	AccountID int64
	// can be negative or positive
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transfer struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	// must be positive
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID                int64
	Email             string
	Fullname          string
	Username          string
	HashedPassword    string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	IsEmailVerified   bool
}

type VerifyEmail struct {
	ID 					int64 	
	Username 		string 
	Email 			string 	
	SecretCode 	string 
	IsUsed 			bool 	
	CreatedAt		time.Time
	ExpiredAt 	time.Time
}
