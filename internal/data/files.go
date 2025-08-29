package data

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
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
	StoragePath    string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	DownloadCount  int       `json:"download_count"`
	Version        int       `json:"-"`
}

type FileModel struct {
	DB *sql.DB
}

func (m FileModel) Insert(file *File) error {
	query := `
		INSERT INTO files (token, owner_id, storage_path, filename)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, last_accessed_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{file.Token, file.OwnerID, file.StoragePath, file.Filename}
	err := m.DB.QueryRowContext(ctx, query, args...).
		Scan(&file.ID, &file.CreatedAt, &file.LastAccessedAt, &file.Version)
	return err
}

func (m FileModel) Update(file *File) error {
	query := `
		UPDATE files
		SET token = $1, owner_id = $2, filename = $3,
		storage_path = $4, created_at = $5,
		last_accessed_at = $6, download_count = $7,
		version = version + 1
		WHERE id = $8 AND version = $9
		RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		file.Token, file.OwnerID,
		file.Filename, file.StoragePath,
		file.CreatedAt, file.LastAccessedAt,
		file.DownloadCount, file.ID, file.Version,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&file.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m FileModel) GetByToken(token string) (*File, error) {
	if len(token) != 8 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, token, owner_id, filename, storage_path, 
		created_at, last_accessed_at, download_count, version
		FROM files WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var file File

	err := m.DB.QueryRowContext(ctx, query, token).
		Scan(&file.ID, &file.Token, &file.OwnerID, &file.Filename, &file.StoragePath, &file.CreatedAt, &file.LastAccessedAt, &file.DownloadCount, &file.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &file, nil
}

func (m FileModel) GetAllForUser(id int64) ([]*File, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, token, owner_id, filename, storage_path, 
		created_at, last_accessed_at, download_count, version
		FROM files WHERE owner_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []*File{}
	for rows.Next() {
		var file File
		err = rows.Scan(
			&file.ID, &file.Token,
			&file.OwnerID, &file.Filename,
			&file.StoragePath, &file.CreatedAt,
			&file.LastAccessedAt, &file.DownloadCount,
			&file.Version,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
