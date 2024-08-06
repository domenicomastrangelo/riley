package sql

import (
	"database/sql"
	"fmt"

	"riley/internal/config"

	_ "github.com/lib/pq"
)

func Connect(config *config.Config) *sql.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s ",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Name,
		config.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	runMigrations(db)

	return db
}

func runMigrations(db *sql.DB) {
	runUserMigration(db)
	runFilesMigration(db)
	runTextsMigration(db)
}

func runUserMigration(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			active BOOLEAN NOT NULL DEFAULT TRUE
		);
	`)
	if err != nil {
		panic(err)
	}
}

func runFilesMigration(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			expires_at TIMESTAMP,
			name VARCHAR(255) NOT NULL,
			hash VARCHAR(255) NOT NULL UNIQUE,
			size BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		panic(err)
	}
}

func runTextsMigration(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS texts (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP,
			expires_at TIMESTAMP,
			data TEXT NOT NULL,
			name VARCHAR(255) NOT NULL,
			size BIGINT NOT NULL,
			hash VARCHAR(255) NOT NULL UNIQUE,
			user_id BIGINT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		panic(err)
	}
}
