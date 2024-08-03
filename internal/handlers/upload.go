package handlers

import (
	"io"
	"net/http"
	"time"

	"riley/internal/auth"
	"riley/internal/models"
)

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)

		_, err := w.Write([]byte("Request body is empty"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	userID, err := auth.GetUserIDFromToken(r.Header.Get("Authorization"), h.Config.TokenSecret)
	if err != nil {
		h.Logger.Error("Error getting user ID from context", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	// Maximum file size is 100MB
	err = r.ParseMultipartForm(100 << 20)
	if err != nil {
		h.Logger.Error("Error parsing multipart form", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Error("Error getting file from form", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}
	defer file.Close()

	var expiresAtTime time.Time
	expiresAt := r.FormValue("expires_at")
	if expiresAt == "" {
		expiresAtTime = time.Now().Add(24 * time.Hour)
	} else {
		expiresAtTime, err = time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			h.Logger.Error("Error parsing expires_at time", "error", err.Error())

			w.WriteHeader(http.StatusBadRequest)

			_, err = w.Write([]byte("Invalid expires_at time"))
			if err != nil {
				h.Logger.Error("Error writing response", "error", err.Error())
			}

			return
		}
	}

	f := models.File{
		Name:      header.Filename,
		Size:      uint64(header.Size),
		UserID:    userID,
		ExpiresAt: expiresAtTime,
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		h.Logger.Error("Error reading file content", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	ff, err := f.CreateFile(&fileContent, h.Config.Storage, h.SQLDatabase)
	if err != nil {
		h.Logger.Error("Error creating file", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(ff.Hash))
	if err != nil {
		h.Logger.Error("Error writing response", "error", err.Error())
	}
}
