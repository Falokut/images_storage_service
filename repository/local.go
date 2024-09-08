package repository

import (
	"context"
	"fmt"
	"github.com/Falokut/images_storage_service/domain"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type LocalImageStorage struct {
	basePath string
}

func NewLocalStorage(baseStoragePath string) *LocalImageStorage {
	return &LocalImageStorage{basePath: baseStoragePath}
}

func (s *LocalImageStorage) Shutdown() {}

func (s *LocalImageStorage) SaveImage(ctx context.Context, img []byte, filename string, relativePath string) error {
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	err := os.MkdirAll(filepath.Dir(relativePath), 0755)
	if err != nil {
		return errors.WithMessage(err, "make directory")
	}

	f, err := os.OpenFile(relativePath, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0660)
	if err != nil {
		return errors.WithMessage(err, "create file")
	}

	_, err = f.Write(img)
	if err != nil {
		return errors.WithMessage(err, "write into file")
	}
	return nil
}

func (s *LocalImageStorage) GetImage(ctx context.Context, filename string, relativePath string) ([]byte, error) {
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))
	image, err := os.ReadFile(relativePath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil, domain.ErrImageNotFound
	case err != nil:
		return nil, errors.WithMessage(err, "read file")
	}
	return image, nil
}

func (s *LocalImageStorage) IsImageExist(ctx context.Context, filename string, relativePath string) (bool, error) {
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))
	_, err := os.Stat(relativePath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, errors.WithMessage(err, "os stat")
	default:
		return true, nil
	}
}

func (s *LocalImageStorage) DeleteImage(ctx context.Context, filename string, relativePath string) error {
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))
	err := os.Remove(relativePath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return domain.ErrImageNotFound
	case err != nil:
		return errors.WithMessage(err, "remove file")
	default:
		return nil
	}
}

func (s *LocalImageStorage) ReplaceImage(ctx context.Context, img []byte, filename string, relativePath string) error {
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	f, err := os.OpenFile(relativePath, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return errors.WithMessage(err, "open file")
	}

	err = f.Truncate(0)
	if err != nil {
		return errors.WithMessage(err, "truncate")
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return errors.WithMessage(err, "seek")
	}
	_, err = f.Write(img)
	if err != nil {
		return errors.WithMessage(err, "write into file")
	}
	return nil
}
