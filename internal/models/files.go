package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"riley/internal/config"
	"riley/internal/storage"

	"github.com/google/uuid"
)

type File struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	ExpiresAt time.Time
	ID        string
	Name      string
	Hash      string
	Size      uint64
	UserID    uint64
}

// CreateFile creates a new file in the database
//
// If the file is created successfully, the file is returned
// If the file is not created successfully, an error is returned
func (f *File) CreateFile(data *[]byte, storageConfig config.StorageConfigInterface, db *sql.DB) (File, error) {
	fileHash, err := createFileHash(data)
	if err != nil {
		return File{}, err
	}

	storage := storage.Storage{
		FileDetails: storage.FileDetails{
			Hash:          fileHash,
			Size:          f.Size,
			FileName:      f.Name,
			FileContent:   data,
			StorageType:   storageConfig.GetStorageType(),
			StorageConfig: storageConfig,
		},
	}

	query := "INSERT INTO files (expires_at, name, hash, size, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at"
	err = db.QueryRow(
		query, f.ExpiresAt, f.Name, fileHash, f.Size, f.UserID,
	).Scan(
		&f.ID, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return File{}, err
	}

	_, err = storage.Upload()
	if err != nil {
		return File{}, err
	}

	file := File{
		Hash:      storage.FileDetails.Hash,
		Size:      storage.FileDetails.Size,
		Name:      storage.FileDetails.FileName,
		ID:        f.ID,
		UserID:    f.UserID,
		ExpiresAt: f.ExpiresAt,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}

	return file, nil
}

func createFileHash(data *[]byte) (string, error) {
	hasher := sha256.New()

	now := fmt.Sprintf("%d", time.Now().UTC().UnixNano())

	rnd := ""

	uuid, err := uuid.NewV7()
	if err != nil {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}

		rnd = hex.EncodeToString(b)
	} else {
		rnd = uuid.String()
	}

	if _, err := hasher.Write([]byte(rnd)); err != nil {
		return "", err
	}

	if _, err := hasher.Write([]byte(now)); err != nil {
		return "", err
	}

	if _, err := hasher.Write(*data); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// GetFileByHash gets a file by the hash
//
// Returns the file if it exists
// Returns an error if the file does not exist
func GetFileByHash(hash string, db *sql.DB) (File, error) {
	file := File{}

	query := "SELECT id, created_at, updated_at, expires_at, name, hash, size, user_id FROM files WHERE hash = $1"
	err := db.QueryRow(query, hash).Scan(&file.ID, &file.CreatedAt, &file.UpdatedAt, &file.ExpiresAt, &file.Name, &file.Hash, &file.Size, &file.UserID)
	if err != nil && err != sql.ErrNoRows {
		return File{}, err
	} else if err == sql.ErrNoRows {
		return File{}, errors.New("file does not exist")
	}

	return file, nil
}

// Delete deletes a file from the database using the ID
//
// Returns an error if the file does not exist
func (f *File) Delete(c config.StorageConfigInterface, db *sql.DB) error {
	query := "DELETE FROM files WHERE hash = $1"
	_, err := db.Exec(query, f.Hash)
	if err != nil {
		return err
	}

	s := storage.Storage{
		FileDetails: storage.FileDetails{
			Hash:          f.Hash,
			FileName:      f.Name,
			StorageType:   c.GetStorageType(),
			StorageConfig: c,
		},
	}

	return s.Delete()
}

// GetFilesByUserID gets all files by the user ID
//
// Returns a list of files if they exist
// Returns an error if the files do not exist
func GetFilesByUserID(id uint64) ([]File, error) {
	files := []File{}

	return files, errors.New("not implemented")
}
