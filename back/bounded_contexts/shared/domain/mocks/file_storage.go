package mocks

import (
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
)

type FileStorageMock struct {
	StoreFunc     func(ctx context.Context, input shareddomain.StoreFileInput) (shareddomain.StoredFile, error)
	FindFunc      func(ctx context.Context, path string) (*shareddomain.StoredFile, error)
	StoreCalls    int
	FindCalls     int
	StoredInputs  []shareddomain.StoreFileInput
	RequestedPath string
}

func (m *FileStorageMock) Store(
	ctx context.Context,
	input shareddomain.StoreFileInput,
) (shareddomain.StoredFile, error) {
	m.StoreCalls++
	m.StoredInputs = append(m.StoredInputs, input)
	if m.StoreFunc == nil {
		return shareddomain.StoredFile{}, nil
	}

	return m.StoreFunc(ctx, input)
}

func (m *FileStorageMock) Find(
	ctx context.Context,
	path string,
) (*shareddomain.StoredFile, error) {
	m.FindCalls++
	m.RequestedPath = path
	if m.FindFunc == nil {
		return nil, nil
	}

	return m.FindFunc(ctx, path)
}
