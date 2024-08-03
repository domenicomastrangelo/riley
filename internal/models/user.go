package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Email     string
	Password  string
	Active    bool
	ID        uint64
}

// UserCheckLogin checks if a user can log in
//
// If the email is invalid, an error is returned
// If the password is invalid, an error is returned
//
// If the user can log in, the user ID is returned
func UserCheckLogin(email string, password string, db *sql.DB) (uint64, error) {
	var (
		encryptedPassword string
		userID            uint64
	)

	query := "SELECT id, password FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(&userID, &encryptedPassword)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	} else if err == sql.ErrNoRows {
		return 0, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// UserExists checks if a user exists in the database
//
// If the email is invalid, an error is returned
// If the user exists, nil is returned
func UserExists(email string, db *sql.DB) error {
	var userID uint64

	query := "SELECT id FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(&userID)
	if err != nil {
		return err
	}

	return nil
}

// UserCreate creates a new user in the database
//
// If the email is invalid, an error is returned
// If the email is already in use, an error is returned
// If the password is invalid, an error is returned
//
// If the user is created successfully, it is returned
func UserCreate(email string, password string, db *sql.DB) (User, error) {
	user := User{}

	password = strings.TrimSpace(password)

	if !UserCreateValidation(email, password) {
		return user, errors.New("invalid email or password, password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	query := "" +
		"INSERT INTO users (email, password)" +
		"VALUES ($1, $2)" +
		"RETURNING id, created_at, updated_at, deleted_at, active"
	err = db.
		QueryRow(query, email, encryptedPassword).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Active)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Delete deletes a user from the database using the ID
//
// If soft is true, the user is soft deleted
// If soft is false, the user is hard deleted
func (u *User) Delete(soft bool, db *sql.DB) error {
	var query string

	if soft {
		query = "" +
			"UPDATE users " +
			"SET active = false " +
			"WHERE id = $1"
	} else {
		query = "DELETE FROM users WHERE id = $1"
	}

	_, err := db.Exec(query, u.ID)

	return err
}

// IsSoftDeleted checks if a user is soft deleted
//
// Returns true if the user is soft deleted
// Returns false if the user is not soft deleted
//
// A user is soft deleted if the Active field is false
func (u *User) IsSoftDeleted(db *sql.DB) bool {
	var active bool

	query := "SELECT active FROM users WHERE id = $1"
	err := db.QueryRow(query, u.ID).Scan(&active)
	if err != nil {
		return false
	}

	return !active
}

// Exists checks if a user exists in the database
//
// Returns true if the user exists
// Returns false if the user does not exist
func (u *User) Exists(db *sql.DB) bool {
	var id uint64

	query := "SELECT id FROM users WHERE id = $1 and active = true"
	err := db.QueryRow(query, u.ID).Scan(&id)
	if err != nil {
		return false
	}

	return id == u.ID
}

// GetFiles returns all the files associated with a user
//
// Returns a slice of Files
func (u *User) GetFiles(db *sql.DB) []File {
	return nil
}

// GetTexts returns all the texts associated with a user
//
// Returns a slice of Texts
func (u *User) GetTexts(db *sql.DB) []Text {
	return nil
}
