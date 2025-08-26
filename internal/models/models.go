package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
	}
	Tokens interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		Insert(token *Token) error
	}
}

func New(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}
