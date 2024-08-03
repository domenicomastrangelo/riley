package models

import (
	"testing"
	"time"

	"riley/internal/sql"

	"riley/internal/config"
)

func TestCreateFile(t *testing.T) {
	fileContent := []byte("test")

	storageConfig := config.LoadTestConfig().Storage

	db := sql.Connect(config.LoadTestConfig())
	user, err := UserCreate("exampleTestCreateFile@example.com", "password123%A%", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	f := File{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Name:      "test.txt",
		Hash:      "test",
		UserID:    user.ID,
	}

	file, err := f.CreateFile(&fileContent, storageConfig, db)
	if err != nil {
		t.Fatalf("CreateFile returned an error: %s", err)
	}

	t.Run("check file properties", func(t *testing.T) {
		if file.ID == "" {
			t.Fatalf("CreateFile did not return a file ID: %s", file.ID)
		}

		if file.Name != f.Name {
			t.Fatalf("expected file name to be %s, got %s", f.Name, file.Name)
		}

		if file.Size != f.Size {
			t.Fatalf("expected file size to be %d, got %d", f.Size, file.Size)
		}

		if file.UserID != f.UserID {
			t.Fatalf("expected user ID to be %d, got %d", f.UserID, file.UserID)
		}

		if file.CreatedAt.IsZero() {
			t.Fatalf("expected created at time to be set, got zero value")
		}

		if file.ExpiresAt != f.ExpiresAt {
			t.Fatalf("expected expires at time to be %v, got %v", f.ExpiresAt, file.ExpiresAt)
		}
	})

	t.Run("delete file", func(t *testing.T) {
		err = file.Delete(storageConfig, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	})
}

func TestGetFileByHash(t *testing.T) {
	fileContent := []byte("test")

	db := sql.Connect(config.LoadTestConfig())

	user, err := UserCreate("exampleTestGetFileByHash@example.com", "password123%A%", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	f := File{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Name:      "test.txt",
		Hash:      "test",
		Size:      uint64(len([]byte("test"))),
		UserID:    user.ID,
	}

	storageConfig := config.LoadTestConfig().Storage

	file, err := f.CreateFile(&fileContent, storageConfig, db)
	if err != nil {
		t.Fatalf("CreateFile returned an error: %s", err)
	}

	t.Run("get file by hash", func(t *testing.T) {
		var file2 File
		file2, err = GetFileByHash(file.Hash, db)
		if err != nil {
			t.Fatalf("GetFileByHash returned an error: %s", err)
		}

		if file.ID != file2.ID {
			t.Fatalf("expected file ID to be %s, got %s", file.ID, file2.ID)
		}

		if file.Name != file2.Name {
			t.Fatalf("expected file name to be %s, got %s", file.Name, file2.Name)
		}

		if file.Size != file2.Size {
			t.Fatalf("expected file size to be %d, got %d", file.Size, file2.Size)
		}

		expectedCreatedAt := file.CreatedAt.UTC().Format(time.RFC3339)
		actualCreatedAt := file2.CreatedAt.UTC().Format(time.RFC3339)

		if expectedCreatedAt != actualCreatedAt {
			t.Fatalf("expected created at time to be %v, got %v", file.CreatedAt, file2.CreatedAt)
		}

		actualExpiresAt := file2.ExpiresAt.UTC().Format(time.RFC3339)
		expectedExpiresAt := file.ExpiresAt.UTC().Format(time.RFC3339)

		if expectedExpiresAt != actualExpiresAt {
			t.Fatalf("expected expires at time to be %v, got %v", file.ExpiresAt, file2.ExpiresAt)
		}

		if file.Hash != file2.Hash {
			t.Fatalf("expected file hash to be %s, got %s", file.Hash, file2.Hash)
		}
	})

	t.Run("delete file", func(t *testing.T) {
		err = file.Delete(storageConfig, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	})
}

func TestGetFilesByUserID(t *testing.T) {
	fileContent := []byte("test")

	db := sql.Connect(config.LoadTestConfig())

	user, err := UserCreate("exampleTestGetFileByHash@example.com", "password123%A%", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	f := File{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Name:      "test2.txt",
		Hash:      "test2",

		UserID: user.ID,
	}

	storageConfig := config.LoadTestConfig().Storage

	file, err := f.CreateFile(&fileContent, storageConfig, db)
	if err != nil {
		t.Fatalf("CreateFile returned an error: %s", err)
	}

	user2, err := UserCreate("exampleTestGetFileByHash2@example.com", "password123%A%", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}

	defer func() {
		err = user2.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	fileContent2 := []byte("test2")

	f2 := File{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Name:      "test3.txt",
		Hash:      "test3",
		Size:      uint64(len([]byte("test3"))),
		UserID:    user2.ID,
	}

	file2, err := f2.CreateFile(&fileContent2, storageConfig, db)
	if err != nil {
		t.Fatalf("CreateFile returned an error: %s", err)
	}

	t.Run("get files by user ID", func(t *testing.T) {
		var files []File
		files, err = GetFilesByUserID(f.UserID)
		if err != nil {
			t.Fatalf("GetFilesByUserID returned an error: %s", err)
		}

		if len(files) != 1 {
			t.Fatalf("expected 1 file, got %d", len(files))
		}

		if files[0].ID != file.ID {
			t.Fatalf("expected file ID to be %s, got %s", file.ID, files[0].ID)
		}

		if files[0].Name != f.Name {
			t.Fatalf("expected file name to be %s, got %s", f.Name, files[0].Name)
		}

		if files[0].Size != f.Size {
			t.Fatalf("expected file size to be %d, got %d", f.Size, files[0].Size)
		}

		if files[0].UserID != f.UserID {
			t.Fatalf("expected user ID to be %d, got %d", f.UserID, files[0].UserID)
		}

		if files[0].CreatedAt != file.CreatedAt {
			t.Fatalf("expected created at time to be %v, got %v", file.CreatedAt, files[0].CreatedAt)
		}

		if files[0].ExpiresAt != file.ExpiresAt {
			t.Fatalf("expected expires at time to be %v, got %v", file.ExpiresAt, files[0].ExpiresAt)
		}

		if files[0].Hash != file.Hash {
			t.Fatalf("expected file hash to be %s, got %s", file.Hash, files[0].Hash)
		}
	})

	t.Run("get files by user ID with no files", func(t *testing.T) {
		err = file.Delete(storageConfig, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}

		err = file2.Delete(storageConfig, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	})
}

func TestDeleteFile(t *testing.T) {
	fileContent := []byte("test4")

	db := sql.Connect(config.LoadTestConfig())

	user, err := UserCreate("exampleTestGetFileByHash@example.com", "password123%A%", db)
	if err != nil {
		t.Fatalf("UserCreate returned an error: %s", err)
	}

	defer func() {
		err = user.Delete(false, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}
	}()

	f := File{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
		Name:      "test.txt",
		Hash:      "test4",
		Size:      uint64(len([]byte("test4"))),
		UserID:    user.ID,
	}

	storageConfig := config.LoadTestConfig().Storage

	file, err := f.CreateFile(&fileContent, storageConfig, db)
	if err != nil {
		t.Fatalf("CreateFile returned an error: %s", err)
	}

	t.Run("delete file", func(t *testing.T) {
		err = file.Delete(storageConfig, db)
		if err != nil {
			t.Fatalf("Delete returned an error: %s", err)
		}

		_, err = GetFileByHash(file.Hash, db)
		if err == nil {
			t.Fatalf("expected GetFileByHash to return an error")
			return
		}
	})
}
