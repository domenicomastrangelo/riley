package models

import (
	"testing"
	"time"

	"riley/internal/config"
	"riley/internal/sql"
)

func TestCreateText(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	textContent := []byte("test")
	textName := "test"

	user, err := UserCreate("testcreatetext@example.com", "password123%%A", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}
	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	expiresAt := time.Now().UTC().Add(time.Hour)

	text, err := CreateText(textName, user.ID, expiresAt, textContent, db)
	if err != nil {
		t.Fatalf("CreateText returned an error: %s", err)
	}

	t.Run("check text properties", func(t *testing.T) {
		if text.ID == "" {
			t.Errorf("expected text ID to be set, got empty string")
		}

		if text.Name != textName {
			t.Errorf("expected text name to be %s, got %s", textName, text.Name)
		}

		expectedSize := uint64(len(textContent))

		if text.Size != expectedSize {
			t.Errorf("expected text size to be %d, got %d", expectedSize, text.Size)
		}

		if text.UserID != user.ID {
			t.Errorf("expected user ID to be %d, got %d", user.ID, text.UserID)
		}

		if text.CreatedAt.IsZero() {
			t.Errorf("expected created at time to be set, got zero value")
		}

		if !text.ExpiresAt.Equal(expiresAt) {
			t.Errorf("expected expires at time to be %v, got %v", expiresAt, text.ExpiresAt)
		}

		if text.Hash == "" {
			t.Errorf("expected text hash to be set, got empty string")
		}
	})

	t.Run("delete text", func(t *testing.T) {
		err = text.Delete(db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	})
}

func TestGetTextByID(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	textContent := []byte("test2")
	textName := "test2"

	user, err := UserCreate("testgettextbyid@example.com", "password123%%A", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}
	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	expiresAt := time.Now().UTC().Add(time.Hour)

	text, err := CreateText(textName, user.ID, expiresAt, textContent, db)
	if err != nil {
		t.Fatalf("CreateText returned an error: %s", err)
	}

	t.Run("check text properties", func(t *testing.T) {
		text2, err := GetTextByHash(text.Hash, db)
		if err != nil {
			t.Fatalf("GetTextByHash returned an error: %s", err)
		}

		if text.ID != text2.ID {
			t.Fatalf("expected text ID to be %s, got %s", text.ID, text2.ID)
		}

		if text.Name != textName {
			t.Fatalf("expected text name to be %s, got %s", textName, text.Name)
		}

		expectedSize := uint64(len(textContent))

		if text.Size != expectedSize {
			t.Fatalf("expected text size to be %d, got %d", expectedSize, text.Size)
		}

		if text.UserID != user.ID {
			t.Fatalf("expected user ID to be %d, got %d", user.ID, text.UserID)
		}

		if text.CreatedAt != text2.CreatedAt {
			t.Fatalf("expected created at time to be %v, got %v", text.CreatedAt, text2.CreatedAt)
		}

		expectedExpiresAt1 := text.ExpiresAt.Format(time.RFC3339)
		expectedExpiresAt2 := text2.ExpiresAt.Format(time.RFC3339)

		if expectedExpiresAt1 != expectedExpiresAt2 {
			t.Fatalf("expected expires at time to be %v, got %v", text.ExpiresAt, text2.ExpiresAt)
		}

		if text.Hash != text2.Hash {
			t.Fatalf("expected text hash to be %s, got %s", text.Hash, text2.Hash)
		}
	})

	t.Run("delete text", func(t *testing.T) {
		err = text.Delete(db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	})
}

func TestGetTextsByUserID(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	textContent := []byte("test3")
	textName := "test3"

	user, err := UserCreate("testgettextsbyuserid@example.com", "password123%%A", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}
	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	otherUser, err := UserCreate("othertestgettextbyuserid@example.com", "password123%%A", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}
	defer func() {
		err = otherUser.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	expiresAt := time.Now().UTC().Add(time.Hour)

	text, err := CreateText(textName, user.ID, expiresAt, textContent, db)
	if err != nil {
		t.Fatalf("CreateText returned an error: %s", err)
	}

	text2, err := CreateText("test4", otherUser.ID, expiresAt, []byte("test4"), db)
	if err != nil {
		t.Fatalf("CreateText returned an error: %s", err)
	}

	t.Run("check text properties", func(t *testing.T) {
		texts, err := GetTextsByUserID(user.ID, db)
		if err != nil {
			t.Fatalf("GetTextsByUserID returned an error: %s", err)
		}

		if len(texts) != 1 {
			t.Fatalf("expected 1 text, got %d", len(texts))
		}

		if text.ID != texts[0].ID {
			t.Fatalf("expected text ID to be %s, got %s", text.ID, texts[0].ID)
		}

		if text.Name != texts[0].Name {
			t.Fatalf("expected text name to be %s, got %s", text.Name, texts[0].Name)
		}

		expectedSize := uint64(len(textContent))

		if expectedSize != texts[0].Size {
			t.Fatalf("expected text size to be %d, got %d", expectedSize, texts[0].Size)
		}

		if text.UserID != texts[0].UserID {
			t.Fatalf("expected user ID to be %d, got %d", text.UserID, texts[0].UserID)
		}

		if text.CreatedAt != texts[0].CreatedAt {
			t.Fatalf("expected created at time to be %v, got %v", text.CreatedAt, texts[0].CreatedAt)
		}

		expectedExpiresAt1 := text.ExpiresAt.Format(time.RFC3339)
		expectedExpiresAt2 := texts[0].ExpiresAt.Format(time.RFC3339)

		if expectedExpiresAt1 != expectedExpiresAt2 {
			t.Fatalf("expected expires at time to be %v, got %v", text.ExpiresAt, texts[0].ExpiresAt)
		}

		if text.Hash != texts[0].Hash {
			t.Fatalf("expected text hash to be %s, got %s", text.Hash, texts[0].Hash)
		}
	})

	t.Run("delete text", func(t *testing.T) {
		err = text.Delete(db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}

		err = text2.Delete(db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	})
}

func TestTextDelete(t *testing.T) {
	db := sql.Connect(config.LoadTestConfig())

	textContent := []byte("test5")
	textName := "test5"

	user, err := UserCreate("testtextdelete@example.com", "password123%%A", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}
	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	expiresAt := time.Now().UTC().Add(time.Hour)

	text, err := CreateText(textName, user.ID, expiresAt, textContent, db)
	if err != nil {
		t.Fatalf("CreateText returned an error: %s", err)
	}

	t.Run("delete text", func(t *testing.T) {
		err = text.Delete(db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}

		_, err = GetTextByHash(text.ID, db)
		if err == nil {
			t.Fatalf("expected GetTextByID to return an error, got nil")
		}
	})
}
