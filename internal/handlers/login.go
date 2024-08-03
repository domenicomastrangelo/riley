package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"riley/internal/auth"
	"riley/internal/models"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	userID, err := models.UserCheckLogin(jsonBody.Email, jsonBody.Password, h.SQLDatabase)
	if err != nil {
		h.Logger.Error("Error checking login", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)

		_, err = w.Write([]byte("Unauthorized"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	token, err := auth.GenerateToken(auth.TokenDefaultExpiryDate(), userID, h.Config.TokenSecret)
	if err != nil {
		h.Logger.Error("Error generating token", "error", err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte("Internal server error"))
		if err != nil {
			h.Logger.Error("Error writing response", "error", err.Error())
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token)

	w.WriteHeader(http.StatusOK)
}
