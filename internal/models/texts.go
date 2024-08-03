package models

import (
	"errors"
	"time"
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
func CreateText(name string, size uint64, userID uint64, expiresAt time.Time, data []byte) (Text, error) {
	text := Text{}

	return text, errors.New("not implemented")
}

// GetTextByID gets a text by the ID
//
// Returns the text if it exists
// Returns an error if the text does not exist
func GetTextByID(id string) (Text, error) {
	text := Text{}

	return text, errors.New("not implemented")
}

// Delete deletes a text from the database using the ID
//
// Returns an error if the text does not exist
func (t *Text) Delete() error {
	return errors.New("not implemented")
}

// GetTextsByUserID gets all texts by the user ID
//
// Returns a list of texts if they exist
// Returns an error if the texts do not exist
func GetTextsByUserID(id uint64) ([]Text, error) {
	texts := []Text{}

	return texts, errors.New("not implemented")
}
