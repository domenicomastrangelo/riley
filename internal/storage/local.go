package storage

import (
	"fmt"
	"os"

	"riley/internal/config"
)

type Local struct {
	FileDetails FileDetails
}

func (l *Local) Upload() (string, error) {
	path := fmt.Sprintf("%s/%s", l.FileDetails.StorageConfig.(*config.StorageConfig).Local.Directory, l.FileDetails.Hash)

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(*l.FileDetails.FileContent)

	return "", err
}

func (l *Local) Delete() error {
	path := fmt.Sprintf("%s/%s", *&l.FileDetails.StorageConfig.(*config.StorageConfig).Local.Directory, l.FileDetails.Hash)

	err := os.Remove(path)

	return err
}

func (l *Local) Exists() error {
	path := fmt.Sprintf("%s/%s", l.FileDetails.StorageConfig.(*config.StorageConfig).Local.Directory, l.FileDetails.Hash)

	_, err := os.Stat(path)

	return err
}

func (l *Local) Download() (*[]byte, error) {
	path := fmt.Sprintf("%s/%s", l.FileDetails.StorageConfig.(*config.StorageConfig).Local.Directory, l.FileDetails.Hash)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content := make([]byte, l.FileDetails.Size)
	_, err = f.Read(content)

	return &content, err
}
