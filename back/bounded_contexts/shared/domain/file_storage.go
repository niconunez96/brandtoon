package shareddomain

import (
	"context"
	"errors"
)

var ErrStoredFileNotFound = errors.New("stored file not found")

type StoreFileInput struct {
	ContentType string
	Data        []byte
	Directory   string
	Name        string
}

type StoredFile struct {
	ContentType string
	Data        []byte
	Path        string
	PublicURL   string
}

type FileStorage interface {
	Store(ctx context.Context, input StoreFileInput) (StoredFile, error)
	Find(ctx context.Context, path string) (*StoredFile, error)
}
