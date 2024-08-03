package storage

/**

import (
	"slices"
	"testing"
	"time"
)

func TestBlobUploadFile(t *testing.T) {
	b := Blob{
		FileDetails: FileDetails{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Hash:      "test",
			Size:      uint64(len([]byte("test"))),
			FileName:  "test.txt",
			FileContent: func() *[]byte {
				content := []byte("test")
				return &content
			}(),
		},
	}

	t.Run("test file upload", func(t *testing.T) {
		_, err := b.Upload()
		if err != nil {
			t.Fatalf("BlobUploadFile returned an error: %s", err)
		}

		err = b.Delete()
		if err != nil {
			t.Fatalf("BlobDeleteFile returned an error: %s", err)
		}
	})
}

func TestBlobExists(t *testing.T) {
	b := Blob{
		FileDetails: FileDetails{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Hash:      "test1",
			Size:      uint64(len([]byte("test1"))),
			FileName:  "test1.txt",
			FileContent: func() *[]byte {
				content := []byte("test1")
				return &content
			}(),
		},
	}

	t.Run("test file exists", func(t *testing.T) {
		_, err := b.Upload()
		if err != nil {
			t.Fatalf("BlobUploadFile returned an error: %s", err)
		}

		err = b.Exists()
		if err != nil {
			t.Fatalf("BlobExists returned an error: %s", err)
		}

		err = b.Delete()
		if err != nil {
			t.Fatalf("BlobDeleteFile returned an error: %s", err)
		}
	})
}

func TestBlobDownloadFile(t *testing.T) {
	b := Blob{
		FileDetails: FileDetails{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Hash:      "test2",
			Size:      uint64(len([]byte("test2"))),
			FileName:  "test2.txt",
			FileContent: func() *[]byte {
				content := []byte("test2")
				return &content
			}(),
		},
	}

	t.Run("test file download", func(t *testing.T) {
		_, err := b.Upload()
		if err != nil {
			t.Fatalf("BlobUploadFile returned an error: %s", err)
		}

		content, err := b.Download()
		if err != nil {
			t.Fatalf("BlobDownloadFile returned an error: %s", err)
		}

		if slices.Compare(*content, *b.FileDetails.FileContent) != 0 {
			t.Fatalf("got different content")
		}

		err = b.Delete()
		if err != nil {
			t.Fatalf("BlobDeleteFile returned an error: %s", err)
		}
	})
}

func TestBlobDeleteFile(t *testing.T) {
	b := Blob{
		FileDetails: FileDetails{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Hash:      "test3",
			Size:      uint64(len([]byte("test3"))),
			FileName:  "test3.txt",
			FileContent: func() *[]byte {
				content := []byte("test3")
				return &content
			}(),
		},
	}

	t.Run("test file delete", func(t *testing.T) {
		_, err := b.Upload()
		if err != nil {
			t.Fatalf("BlobUploadFile returned an error: %s", err)
		}

		err = b.Delete()
		if err != nil {
			t.Fatalf("BlobDeleteFile returned an error: %s", err)
		}

		err = b.Exists()
		if err != nil {
			t.Fatalf("BlobExists returned an error: %s", err)
		}
	})
}

*/
