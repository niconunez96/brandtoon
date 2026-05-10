package sharedstorage

import (
	shareddomain "brandtoonapi/bounded_contexts/shared/domain"
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const defaultPublicFileURLPrefix = "/files"

type LocalFileStorage struct {
	publicBaseURL string
	publicPrefix  string
	rootPath      string
}

func NewLocalFileStorage(rootPath string, publicBaseURL string, publicPrefix string) (*LocalFileStorage, error) {
	trimmedRoot := strings.TrimSpace(rootPath)
	if trimmedRoot == "" {
		return nil, errors.New("local file storage root path is required")
	}

	trimmedBaseURL := strings.TrimSpace(publicBaseURL)
	if trimmedBaseURL == "" {
		return nil, errors.New("public base url is required")
	}

	trimmedPrefix := strings.TrimSpace(publicPrefix)
	if trimmedPrefix == "" {
		trimmedPrefix = defaultPublicFileURLPrefix
	}
	if !strings.HasPrefix(trimmedPrefix, "/") {
		trimmedPrefix = "/" + trimmedPrefix
	}

	return &LocalFileStorage{
		publicBaseURL: strings.TrimRight(trimmedBaseURL, "/"),
		publicPrefix:  strings.TrimRight(trimmedPrefix, "/"),
		rootPath:      trimmedRoot,
	}, nil
}

func (s *LocalFileStorage) Store(
	ctx context.Context,
	input shareddomain.StoreFileInput,
) (shareddomain.StoredFile, error) {
	select {
	case <-ctx.Done():
		return shareddomain.StoredFile{}, ctx.Err()
	default:
	}

	relDir, err := sanitizeRelativePath(input.Directory)
	if err != nil {
		return shareddomain.StoredFile{}, err
	}

	fileName, err := sanitizeFileName(input.Name, input.ContentType)
	if err != nil {
		return shareddomain.StoredFile{}, err
	}

	relPath := path.Join(relDir, fileName)
	absolutePath := filepath.Join(s.rootPath, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return shareddomain.StoredFile{}, fmt.Errorf("create storage directory: %w", err)
	}

	if err := os.WriteFile(absolutePath, input.Data, 0o644); err != nil {
		return shareddomain.StoredFile{}, fmt.Errorf("write stored file: %w", err)
	}

	return shareddomain.StoredFile{
		ContentType: input.ContentType,
		Data:        append([]byte(nil), input.Data...),
		Path:        relPath,
		PublicURL:   buildPublicURL(s.publicBaseURL, s.publicPrefix, relPath),
	}, nil
}

func (s *LocalFileStorage) Find(
	ctx context.Context,
	storedPath string,
) (*shareddomain.StoredFile, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	relPath, err := sanitizeRelativePath(storedPath)
	if err != nil {
		return nil, err
	}

	absPath := filepath.Join(s.rootPath, filepath.FromSlash(relPath))
	data, err := os.ReadFile(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, shareddomain.ErrStoredFileNotFound
		}

		return nil, fmt.Errorf("read stored file: %w", err)
	}

	contentType := mime.TypeByExtension(filepath.Ext(absPath))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	return &shareddomain.StoredFile{
		ContentType: contentType,
		Data:        data,
		Path:        relPath,
		PublicURL:   buildPublicURL(s.publicBaseURL, s.publicPrefix, relPath),
	}, nil
}

func NewPublicFileHandler(storage shareddomain.FileStorage, publicPrefix string) http.Handler {
	trimmedPrefix := strings.TrimRight(strings.TrimSpace(publicPrefix), "/")
	if trimmedPrefix == "" {
		trimmedPrefix = defaultPublicFileURLPrefix
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet, http.MethodHead:
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		storedPath := strings.TrimPrefix(request.URL.Path, trimmedPrefix)
		storedPath = strings.TrimPrefix(storedPath, "/")
		if storedPath == "" {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		storedFile, err := storage.Find(request.Context(), storedPath)
		if err != nil {
			if errors.Is(err, shareddomain.ErrStoredFileNotFound) {
				writer.WriteHeader(http.StatusNotFound)
				return
			}

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if storedFile == nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		if storedFile.ContentType != "" {
			writer.Header().Set("Content-Type", storedFile.ContentType)
		}
		writer.Header().Set("Cache-Control", "public, max-age=31536000")
		if request.Method == http.MethodHead {
			writer.WriteHeader(http.StatusOK)
			return
		}

		_, _ = writer.Write(storedFile.Data)
	})
}

func sanitizeRelativePath(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", nil
	}

	if path.IsAbs(trimmed) {
		return "", errors.New("absolute paths are not allowed")
	}

	cleaned := path.Clean(trimmed)
	if cleaned == "." {
		return "", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", errors.New("path traversal is not allowed")
	}

	return strings.TrimPrefix(cleaned, "./"), nil
}

func sanitizeFileName(name string, contentType string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", errors.New("file name is required")
	}

	if strings.Contains(trimmed, "/") {
		return "", errors.New("file name must not contain path separators")
	}

	if path.Ext(trimmed) != "" {
		return trimmed, nil
	}

	extensions, _ := mime.ExtensionsByType(contentType)
	if len(extensions) > 0 {
		return trimmed + extensions[0], nil
	}

	return trimmed, nil
}

func buildPublicURL(baseURL string, prefix string, storedPath string) string {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return strings.TrimRight(baseURL, "/") + path.Join(prefix, storedPath)
	}

	parsed.Path = path.Join(parsed.Path, prefix, storedPath)
	return parsed.String()
}
