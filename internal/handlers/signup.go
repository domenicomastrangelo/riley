package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"riley/internal/auth"
	"riley/internal/models"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	jsonBody := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    `json:"email"`,
		Password: `json:"password"`,
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Error("Error reading body", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		h.Logger.Error("Error marshalling body into JSON", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	err = models.UserExists(jsonBody.Email, h.SQLDatabase)
	if err != nil {
		w.WriteHeader(http.StatusConflict)

		_, err = w.Write([]byte("User already exists"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	user, err := models.UserCreate(jsonBody.Email, jsonBody.Password, h.SQLDatabase)
	if err != nil {
		h.Logger.Error("Error creating user", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	w.WriteHeader(http.StatusCreated)

	token, err := auth.GenerateToken(
		auth.TokenDefaultExpiryDate(),
		user.ID,
		h.Config.TokenSecret,
	)
	if err != nil {
		h.Logger.Error("Error generating token", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	_, err = w.Write([]byte(token))
	if err != nil {
		h.Logger.Error("Error writing response", "error", err.Error())
	}
}
