package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrEditConflict     = errors.New("edit conflict")
	ErrPermissionDenied = errors.New("permission denied")
)

type Models struct {
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		GetForToken(scope, tokenPlainText string) (*User, error)
	}
	Tokens interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		Insert(token *Token) error
	}

	Files interface {
		Insert(file *File) error
		GetByToken(token string) (*File, error)
		Update(file *File) error
		GetAllForUser(id int64) ([]*File, error)
		Delete(token string, ownerID int64) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Files:  FileModel{DB: db},
	}
}
