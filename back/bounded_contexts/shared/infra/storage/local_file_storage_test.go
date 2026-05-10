package sharedstorage

import (
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalFileStorageStoresAndFindsFiles(t *testing.T) {
	t.Parallel()

	storage, err := NewLocalFileStorage(t.TempDir(), "http://127.0.0.1:8888", "/files")
	if err != nil {
		t.Fatalf("expected storage, got error %v", err)
	}

	storedFile, err := storage.Store(context.Background(), shareddomain.StoreFileInput{
		ContentType: "image/png",
		Data:        []byte("png-bytes"),
		Directory:   "avatars/avatar-v7/options",
		Name:        "option-v7",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if storedFile.Path != "avatars/avatar-v7/options/option-v7.png" {
		t.Fatalf("unexpected path %q", storedFile.Path)
	}

	if storedFile.PublicURL != "http://127.0.0.1:8888/files/avatars/avatar-v7/options/option-v7.png" {
		t.Fatalf("unexpected public url %q", storedFile.PublicURL)
	}

	foundFile, err := storage.Find(context.Background(), storedFile.Path)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if string(foundFile.Data) != "png-bytes" {
		t.Fatalf("unexpected file bytes %q", string(foundFile.Data))
	}

	if foundFile.ContentType != "image/png" {
		t.Fatalf("unexpected content type %q", foundFile.ContentType)
	}
}

func TestLocalFileStorageRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	storage, err := NewLocalFileStorage(t.TempDir(), "http://127.0.0.1:8888", "/files")
	if err != nil {
		t.Fatalf("expected storage, got error %v", err)
	}

	_, err = storage.Store(context.Background(), shareddomain.StoreFileInput{
		ContentType: "image/png",
		Data:        []byte("png-bytes"),
		Directory:   "../escape",
		Name:        "option-v7",
	})
	if err == nil {
		t.Fatalf("expected path traversal error")
	}
}

func TestPublicFileHandlerServesStoredFiles(t *testing.T) {
	t.Parallel()

	storage, err := NewLocalFileStorage(t.TempDir(), "http://127.0.0.1:8888", "/files")
	if err != nil {
		t.Fatalf("expected storage, got error %v", err)
	}

	storedFile, err := storage.Store(context.Background(), shareddomain.StoreFileInput{
		ContentType: "image/png",
		Data:        []byte("png-bytes"),
		Directory:   "avatars/avatar-v7/options",
		Name:        "option-v7",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	handler := NewPublicFileHandler(storage, "/files")
	request := httptest.NewRequest(http.MethodGet, "/files/"+storedFile.Path, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("unexpected content type %q", recorder.Header().Get("Content-Type"))
	}

	if recorder.Body.String() != "png-bytes" {
		t.Fatalf("unexpected body %q", recorder.Body.String())
	}
}

func TestPublicFileHandlerReturnsNotFoundForMissingFiles(t *testing.T) {
	t.Parallel()

	handler := NewPublicFileHandler(&failingFileStorage{}, "/files")
	request := httptest.NewRequest(http.MethodGet, "/files/missing.png", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recorder.Code)
	}
}

type failingFileStorage struct{}

func (f *failingFileStorage) Store(
	ctx context.Context,
	input shareddomain.StoreFileInput,
) (shareddomain.StoredFile, error) {
	return shareddomain.StoredFile{}, errors.New("not implemented")
}

func (f *failingFileStorage) Find(
	ctx context.Context,
	path string,
) (*shareddomain.StoredFile, error) {
	return nil, shareddomain.ErrStoredFileNotFound
}
