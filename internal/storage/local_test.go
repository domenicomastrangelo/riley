package storage

import (
	"testing"

	"riley/internal/config"
)

func TestLocalUploadFile(t *testing.T) {
	l := Local{
		FileDetails: FileDetails{
			Hash:     "test",
			Size:     uint64(len([]byte("test"))),
			FileName: "test.txt",
			FileContent: func() *[]byte {
				content := []byte("test")
				return &content
			}(),
			StorageConfig: &config.StorageConfig{
				StorageType: STORAGE_TYPE_LOCAL,
				Local: config.LocalConfig{
					Directory: "/tmp",
				},
			},
		},
	}

	t.Run("test file upload", func(t *testing.T) {
		_, err := l.Upload()
		if err != nil {
			t.Fatalf("LocalUploadFile returned an error: %s", err)
		}

		err = l.Delete()
		if err != nil {
			t.Fatalf("LocalDeleteFile returned an error: %s", err)
		}
	})
}

func TestLocalExists(t *testing.T) {
	l := Local{
		FileDetails: FileDetails{
			Hash:     "test1",
			Size:     uint64(len([]byte("test1"))),
			FileName: "test1.txt",
			FileContent: func() *[]byte {
				content := []byte("test1")
				return &content
			}(),
			StorageConfig: &config.StorageConfig{
				StorageType: STORAGE_TYPE_LOCAL,
				Local: config.LocalConfig{
					Directory: "/tmp",
				},
			},
		},
	}

	t.Run("test file exists", func(t *testing.T) {
		_, err := l.Upload()
		if err != nil {
			t.Fatalf("LocalUploadFile returned an error: %s", err)
		}

		err = l.Exists()
		if err != nil {
			t.Fatalf("LocalExists returned an error: %s", err)
		}

		err = l.Delete()
		if err != nil {
			t.Fatalf("LocalDeleteFile returned an error: %s", err)
		}
	})
}

func TestLocalDownload(t *testing.T) {
	l := Local{
		FileDetails: FileDetails{
			Hash:     "test2",
			Size:     uint64(len([]byte("test2"))),
			FileName: "test2.txt",
			FileContent: func() *[]byte {
				content := []byte("test2")
				return &content
			}(),
			StorageConfig: &config.StorageConfig{
				StorageType: STORAGE_TYPE_LOCAL,
				Local: config.LocalConfig{
					Directory: "/tmp",
        },
      },
		},
	}

	t.Run("test file download", func(t *testing.T) {
		_, err := l.Upload()
		if err != nil {
			t.Fatalf("LocalUploadFile returned an error: %s", err)
		}

		content, err := l.Download()
		if err != nil {
			t.Fatalf("LocalDownload returned an error: %s", err)
		}

		if string(*content) != "test2" {
			t.Fatalf("got different content")
		}

		err = l.Delete()
		if err != nil {
			t.Fatalf("LocalDeleteFile returned an error: %s", err)
		}
	})
}

func TestLocalDeleteFile(t *testing.T) {
	l := Local{
		FileDetails: FileDetails{
			Hash:     "test3",
			Size:     uint64(len([]byte("test3"))),
			FileName: "test3.txt",
			FileContent: func() *[]byte {
				content := []byte("test3")
				return &content
			}(),
			StorageConfig: &config.StorageConfig{
				StorageType: STORAGE_TYPE_LOCAL,
				Local: config.LocalConfig{
					Directory: "/tmp",
				},
			},
		},
	}

	t.Run("test file delete", func(t *testing.T) {
		_, err := l.Upload()
		if err != nil {
			t.Fatalf("LocalUploadFile returned an error: %s", err)
		}

		err = l.Delete()
		if err != nil {
			t.Fatalf("LocalDeleteFile returned an error: %s", err)
		}

		err = l.Exists()
		if err == nil {
			t.Fatalf("LocalExists returned nil, expected error")
		}
	})
}
