package storage

import (
	"errors"
	"time"

	"riley/internal/config"
)

type StorageInterface interface {
	Upload() (string, error)
	Delete() error
	Exists() error
	Download() (*[]byte, error)
}

type Storage struct {
	FileDetails FileDetails
}

const (
	STORAGE_TYPE_LOCAL = "local"
	STORAGE_TYPE_BLOB  = "blob"
)

type FileDetails struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StorageConfig config.StorageConfigInterface
	StorageType   string
	FileContent   *[]byte
	Hash          string
	FileName      string
	Storage       string
	Size          uint64
}

func (s *Storage) Upload() (string, error) {
	switch s.FileDetails.StorageType {
	case STORAGE_TYPE_LOCAL:
		l := Local{
			FileDetails: s.FileDetails,
		}
		return l.Upload()
	case STORAGE_TYPE_BLOB:
		b := Blob{
			FileDetails: s.FileDetails,
		}
		return b.Upload()
	}

	return "", errors.New("storage type not implemented")
}

func (s *Storage) Delete() error {
	switch s.FileDetails.StorageType {
	case STORAGE_TYPE_LOCAL:
		l := Local{
			FileDetails: s.FileDetails,
		}
		return l.Delete()
	case STORAGE_TYPE_BLOB:
		b := Blob{
			FileDetails: s.FileDetails,
		}
		return b.Delete()
	}

	return nil
}

func (s *Storage) Exists() error {
	return nil
}

func (s *Storage) Download() ([]byte, error) {
	return nil, nil
}
