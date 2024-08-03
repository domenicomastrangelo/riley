package storage

import (
	"errors"
)

type Blob struct {
	FileDetails FileDetails
}

func (b *Blob) Upload() (string, error) {
	return "", errors.New("not implemented")
}

func (b *Blob) Delete() error {
	return errors.New("not implemented")
}

func (b *Blob) Exists() error {
	return errors.New("not implemented")
}

func (b *Blob) Download() (*[]byte, error) {
	return nil, errors.New("not implemented")
}
