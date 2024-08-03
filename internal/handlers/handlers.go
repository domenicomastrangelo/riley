package handlers

import (
	"database/sql"
	"log/slog"

	"riley/internal/config"
)

type Handler struct {
	SQLDatabase *sql.DB
	Config      *config.Config
	Logger      *slog.Logger
}
