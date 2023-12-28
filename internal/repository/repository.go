package repository

import (
	"context"
	"errors"
)

var (
	ErrNotExist = errors.New("file doesn't exist")
)

//go:generate mockgen -source=repository.go -destination=mocks/imageStorage.go
type ImageStorage interface {
	SaveImage(ctx context.Context, img []byte, filename string, category string) error
	GetImage(ctx context.Context, imageID string, category string) ([]byte, error)
	IsImageExist(ctx context.Context, imageID string, category string) bool
	DeleteImage(ctx context.Context, imageID string, category string) error
	RewriteImage(ctx context.Context, img []byte, filename string, category string) error
}
