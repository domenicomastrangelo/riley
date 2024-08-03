package handlers

import (
	"bytes"
	"encoding/hex"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"riley/internal/auth"
	"riley/internal/config"
	"riley/internal/models"
	"riley/internal/sql"
)

func TestUpload(t *testing.T) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create fake file in memory
	fileWriter, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = fileWriter.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/upload", &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	h := createHandler()

	// Create user
	user, err := models.UserCreate("testupload@example.com", "password123%A%", h.SQLDatabase)
	if err != nil {
		t.Fatal(err)
	}

	token, err := auth.GenerateToken(time.Now().Add(1*time.Hour), user.ID, h.Config.TokenSecret)
	if err != nil {
		t.Fatal(err)
	}

	// Set Authorization header
	req.Header.Set("Authorization", token)

	// Set content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(h.Upload)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body is a valid SHA256 hash
	if !isValidSHA256(rr.Body.String()) {
		t.Errorf("handler returned invalid SHA256 hash: got %v", rr.Body.String())
	}

	file := models.File{
		Hash: rr.Body.String(),
	}

	// Delete file
	err = file.Delete(config.LoadTestConfig().Storage, h.SQLDatabase)
	if err != nil {
		t.Fatal(err)
	}

	// Delete user
	err = user.Delete(false, h.SQLDatabase)
	if err != nil {
		t.Fatal(err)
	}
}

func isValidSHA256(s string) bool {
	// Check if the string is 64 characters long
	if len(s) != 64 {
		return false
	}

	// Try to decode the string from hexadecimal
	_, err := hex.DecodeString(s)
	return err == nil
}

func createHandler() *Handler {
	db := sql.Connect(config.LoadTestConfig())
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	h := Handler{
		SQLDatabase: db,
		Config:      config.LoadTestConfig(),
		Logger:      logger,
	}

	return &h
}
