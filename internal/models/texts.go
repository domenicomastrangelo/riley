package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Text struct {
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

// CreateText creates a new text in the database
//
// If the text is created successfully, the text is returned
// If the text is not created successfully, an error is returned
func CreateText(name string, userID uint64, expiresAt time.Time, data []byte, db *sql.DB) (Text, error) {
	text := Text{}

	if validateText(name, userID, expiresAt, data) != nil {
		return text, errors.New("invalid text")
	}

	size := uint64(len(data))

	hash, err := createTextHash(&data)
	if err != nil {
		return text, err
	}

	query := "INSERT INTO texts (expires_at, name, hash, size, user_id, data) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at"
	err = db.QueryRow(
		query, expiresAt, name, hash, size, userID, data,
	).Scan(
		&text.ID, &text.CreatedAt, &text.UpdatedAt,
	)
	if err != nil {
		return text, err
	}

	text.Size = size
	text.Name = name
	text.Hash = hash
	text.UserID = userID
	text.ExpiresAt = expiresAt

	return text, nil
}

func createTextHash(data *[]byte) (string, error) {
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

func validateText(name string, userID uint64, expiresAt time.Time, data []byte) error {
	if len(name) == 0 {
		return errors.New("name is required")
	}

	if userID == 0 {
		return errors.New("userID is required")
	}

	if expiresAt.IsZero() {
		return errors.New("expiresAt is required")
	}

	if len(data) == 0 {
		return errors.New("data is required")
	}

	return nil
}

// GetTextByID gets a text by the ID
//
// Returns the text if it exists
// Returns an error if the text does not exist
func GetTextByHash(hash string, db *sql.DB) (Text, error) {
	text := Text{}

	query := "SELECT id, created_at, updated_at, expires_at, name, hash, size, user_id FROM texts WHERE hash = $1"
	err := db.QueryRow(query, hash).Scan(&text.ID, &text.CreatedAt, &text.UpdatedAt, &text.ExpiresAt, &text.Name, &text.Hash, &text.Size, &text.UserID)
	if err != nil && err != sql.ErrNoRows {
		return Text{}, err
	} else if err == sql.ErrNoRows {
		return Text{}, errors.New("text does not exist")
	}

	return text, nil
}

// Delete deletes a text from the database using the ID
//
// Returns an error if the text does not exist
func (t *Text) Delete(db *sql.DB) error {
	query := "DELETE FROM texts WHERE id = $1"
	_, err := db.Exec(query, t.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetTextsByUserID gets all texts by the user ID
//
// Returns a list of texts if they exist
// Returns an error if the texts do not exist
func GetTextsByUserID(id uint64, db *sql.DB) ([]Text, error) {
	texts := []Text{}

	query := "SELECT id, created_at, updated_at, expires_at, name, hash, size, user_id FROM texts WHERE user_id = $1"

	rows, err := db.Query(query, id)
	if err != nil && err != sql.ErrNoRows {
		return []Text{}, err
	}

	for rows.Next() {
		var text Text
		err := rows.Scan(&text.ID, &text.CreatedAt, &text.UpdatedAt, &text.ExpiresAt, &text.Name, &text.Hash, &text.Size, &text.UserID)
		if err != nil {
			return []Text{}, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}
