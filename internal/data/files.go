package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NewFileToken() (string, error) {
	token := make([]byte, 8)
	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(62))
		if err != nil {
			return "", err
		}
		token[i] = base62Chars[num.Int64()]
	}
	return string(token), nil
}

type File struct {
	ID             int64     `json:"id"`
	Token          string    `json:"token"`
	OwnerID        *int64    `json:"owner_id"`
	Filename       string    `json:"filename"`
	StoragePath    string    `json:"storage_path"`
	CreatedAt      time.Time `json:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	DownloadCount  int64     `json:"download_count"`
}

type FileModel struct {
	DB *sql.DB
}

func (m FileModel) Insert(file *File) error {
	query := `
		INSERT INTO files (token, owner_id, storage_path, filename)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, last_accessed_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{file.Token, file.OwnerID, file.StoragePath, file.Filename}
	err := m.DB.QueryRowContext(ctx, query, args...).
		Scan(&file.ID, &file.CreatedAt, &file.LastAccessedAt)
	return err
}
